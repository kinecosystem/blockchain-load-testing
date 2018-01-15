package account

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"

	txerrors "github.com/kinfoundation/stellar-benchmark/src/errors"
)

const (
	fundTimeout = 10 * time.Second
	getTimeout  = fundTimeout

	submitTimeout       = 10 * time.Second
	retryFailedTxAmount = 10
)

func Create(horizonAddr string, funder *keypair.Full, accountsNum int, fundAmount string, logger log.Logger) ([]keypair.KP, error) {
	level.Info(logger).Log("msg", "creating accounts", "accounts_num", accountsNum)

	ops := make([]build.TransactionMutator, 0)

	keypairs := make([]keypair.KP, accountsNum)
	for i := 0; i < accountsNum; i++ {
		level.Info(logger).Log("msg", "adding create account operation", "accounts_num", i)

		kp, err := keypair.Random()
		if err != nil {
			level.Error(logger).Log("msg", err, "account_index", i, "seed", kp.Seed())
			return nil, err
		}
		keypairs[i] = kp

		ops = append(ops, build.CreateAccount(
			build.Destination{AddressOrSeed: kp.Address()},
			build.NativeAmount{Amount: fundAmount}))
	}

	level.Info(logger).Log("msg", "submitting create account transaction")

	err := submitWithRetry(horizonAddr, ops, funder.Seed(), logger)
	if err != nil {
		txerrors.GetTxErrorResultCodes(err, logger)
		return nil, err
	}

	for i, kp := range keypairs {
		level.Info(logger).Log(
			"msg", "new account created",
			"account_index", i,
			"address", kp.Address()[:5],
			"seed", kp.(*keypair.Full).Seed())
	}

	return keypairs, nil
}

func Get(horizonAddr, address string, logger log.Logger) (*horizon.Account, error) {
	l := log.With(logger, "address", address[:5])

	client := horizon.Client{
		URL:  horizonAddr,
		HTTP: &http.Client{Timeout: getTimeout},
	}
	account, err := client.LoadAccount(address)
	if err != nil {
		level.Error(l).Log("msg", err)
		return nil, err
	}

	return &account, nil
}

func submitWithRetry(horizonAddr string, ops []build.TransactionMutator, seed string, logger log.Logger) error {
	var err error
	for i := 0; i < retryFailedTxAmount; i++ {
		level.Info(logger).Log("retry_index", i, "msg", "submitting transaction")

		client := horizon.Client{
			URL:  horizonAddr,
			HTTP: &http.Client{Timeout: submitTimeout}}

		fullOps := append(
			[]build.TransactionMutator{
				build.SourceAccount{
					AddressOrSeed: seed},
				build.TestNetwork,
				build.AutoSequence{
					SequenceProvider: &horizon.Client{
						URL:  horizonAddr,
						HTTP: client.HTTP,
					},
				},
			},
			ops...)

		txBuilder := build.Transaction(fullOps...)
		txEnv := txBuilder.Sign(seed)
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

		txerrors.GetTxErrorResultCodes(err, logger)

		time.Sleep(5 * time.Second)
	}

	return err
}
