package main

import (
	"math"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

const (
	maxOpsInTx = 100

	fundTimeout = 10 * time.Second
	getTimeout  = fundTimeout

	submitTimeout       = 20 * time.Second
	retryFailedTxAmount = 10
)

func Create(horizonAddr string, network build.Network, funder *keypair.Full, accountsNum int, fundAmount string, logger log.Logger) ([]keypair.KP, error) {
	level.Info(logger).Log("msg", "creating accounts", "accounts_num", accountsNum)

	keypairs := make([]keypair.KP, 0, accountsNum)

	for batchIndex := 0; batchIndex <= (cap(keypairs)+1)/(maxOpsInTx+1); batchIndex++ {
		batch := keypairs[batchIndex*maxOpsInTx : int(math.Min(float64((batchIndex+1)*maxOpsInTx), float64(cap(keypairs))))]

		ops := make([]build.TransactionMutator, 0)
		for i := 0; i < len(batch); i++ {
			level.Info(logger).Log("msg", "adding create account operation", "account_index", batchIndex*maxOpsInTx+i)

			kp, err := keypair.Random()
			if err != nil {
				level.Error(logger).Log("msg", err, "account_index", batchIndex*maxOpsInTx+i, "seed", kp.Seed())
				return nil, err
			}
			keypairs = append(keypairs, kp)

			ops = append(ops, build.CreateAccount(
				build.Destination{AddressOrSeed: kp.Address()},
				build.NativeAmount{Amount: fundAmount}))
		}

		level.Info(logger).Log("msg", "submitting create account transaction")

		err := submitWithRetry(horizonAddr, network, ops, funder.Seed(), logger)
		if err != nil {
			GetTxErrorResultCodes(err, logger)
			return nil, err
		}

		for i, kp := range batch {
			level.Info(logger).Log(
				"msg", "new account created",
				"account_index", batchIndex*maxOpsInTx+i,
				"address", kp.Address(),
				"seed", kp.(*keypair.Full).Seed())
		}
	}

	return keypairs, nil
}

func submitWithRetry(horizonAddr string, network build.Network, ops []build.TransactionMutator, seed string, logger log.Logger) error {
	var err error
	for i := 0; i < retryFailedTxAmount; i++ {
		level.Info(logger).Log("retry_index", i, "msg", "submitting transaction")

		client := horizon.Client{
			URL:  horizonAddr,
			HTTP: &http.Client{Timeout: submitTimeout}}

		fullOps := append(
			[]build.TransactionMutator{
				build.SourceAccount{AddressOrSeed: seed},
				network,
				build.AutoSequence{
					SequenceProvider: &horizon.Client{
						URL:  horizonAddr,
						HTTP: client.HTTP,
					},
				},
			},
			ops...,
		)

		txBuilder, err := build.Transaction(fullOps...)
		if err != nil {
			level.Error(logger).Log("msg", err)
			return err
		}

		txEnv, err := txBuilder.Sign(seed)
		if err != nil {
			level.Error(logger).Log("msg", err)
			return err
		}

		var txEnvB64 string
		txEnvB64, err = txEnv.Base64()
		if err != nil {
			level.Error(logger).Log("msg", err)
			return err
		}

		_, err = client.SubmitTransaction(txEnvB64)
		if err == nil {
			return nil
		}

		GetTxErrorResultCodes(err, logger)

		time.Sleep(5 * time.Second)
	}

	return err
}
