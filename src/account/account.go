package account

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"

	"github.com/kinfoundation/stellar-benchmark/src/transaction"
	txerrors "github.com/kinfoundation/stellar-benchmark/src/transaction/errors"
)

const (
	fundTimeout = 10 * time.Second
	getTimeout  = fundTimeout
)

func Create(horizonAddr string, funder *keypair.Full, accountsNum int, fundAmount string, logger log.Logger) ([]keypair.KP, error) {
	level.Info(logger).Log("msg", "creating accounts", "accounts_num", accountsNum)

	ops := make([]build.TransactionMutator, 0, accountsNum+3)
	ops = append(
		ops,
		build.SourceAccount{AddressOrSeed: funder.Address()},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient})

	keypairs := make([]keypair.KP, accountsNum)
	for i := 0; i < accountsNum; i++ {
		level.Info(logger).Log("msg", "adding create account operation", "accounts_num", i)

		kp, err := keypair.Random()
		if err != nil {
			level.Error(logger).Log("msg", err, "account_index", i, "seed", kp.Seed())
			return nil, err
		}
		keypairs[i] = kp

		ops = append(
			ops,
			build.CreateAccount(
				build.Destination{AddressOrSeed: kp.Address()},
				build.NativeAmount{Amount: fundAmount}))
	}

	level.Info(logger).Log("msg", "submitting create account transaction")

	txBuilder := build.Transaction(ops...)
	err := transaction.SubmitWithRetry(horizonAddr, txBuilder, funder.Seed(), logger)
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

func Fund(horizonAddr string, kp keypair.KP, logger log.Logger) error {
	l := log.With(logger, "address", kp.Address()[:5])

	level.Info(l).Log("msg", "sending funding request", "address", kp.Address()[:5])

	client := http.Client{Timeout: fundTimeout}

	res, err := client.Get(fmt.Sprintf("%s/friendbot?addr=%s", horizonAddr, kp.Address()))
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			level.Error(l).Log("msg", err)
		}
	}()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	level.Debug(l).Log("msg", string(data))

	if res.StatusCode == http.StatusBadRequest {
		return errors.New("funding failure")
	}

	level.Info(l).Log("msg", "funding success")

	return nil
}
