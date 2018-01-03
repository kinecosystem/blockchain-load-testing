package transaction

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"

	"github.com/kinfoundation/stellar-benchmark/src/transaction/errors"
)

const (
	submitTimeout       = 10 * time.Second
	retryFailedTxAmount = 10
)

func Transfer(horizonAddr string, from, to keypair.KP, amount string, logger log.Logger) error {
	txBuilder := build.Transaction(
		build.SourceAccount{AddressOrSeed: from.Address()},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},

		build.Payment(
			build.Destination{AddressOrSeed: to.Address()},
			build.NativeAmount{Amount: amount},
		),
	)

	l := log.With(
		logger,
		"from", from.Address()[:5],
		"to", to.Address()[:5],
		"amount", amount)

	err := SubmitWithRetry(horizonAddr, txBuilder, from.(*keypair.Full).Seed(), l)
	if err != nil {
		errors.GetTxErrorResultCodes(err, logger)
		return err
	}
	return nil
}

func SubmitWithRetry(horizonAddr string, txBuilder *build.TransactionBuilder, seed string, logger log.Logger) error {
	client := horizon.Client{
		URL:  horizonAddr,
		HTTP: &http.Client{Timeout: submitTimeout},
	}

	txEnv := txBuilder.Sign(seed)
	txEnvB64, err := txEnv.Base64()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}

	for i := 0; i < retryFailedTxAmount; i++ {
		level.Info(logger).Log("retry_index", i, "msg", "submitting transaction")

		_, err = client.SubmitTransaction(txEnvB64)
		if err == nil {
			return nil
		}

		errors.GetTxErrorResultCodes(err, logger)
	}

	return err
}
