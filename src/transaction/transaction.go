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
	funderSeed = "SB46AOUVB73DOLQVJWBW676HCA6IA4WISBUGLMXSNRW4P3JOLM3MO2RV"

	accountsNum = 100

	fundAmount     = "30"
	transferAmount = "0.01"

	createAccountTimeout = 10 * time.Second
	transferTimeout      = createAccountTimeout

	retryFailedTxAmount = 3
)

func Transfer(from, to keypair.KP, amount string, logger log.Logger) error {
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

	level.Info(l).Log("msg", "submitting transaction")

	err := SubmitWithRetry(txBuilder, from.(*keypair.Full).Seed(), l)
	if err != nil {
		errors.GetTxErrorResultCodes(err, logger)
		return err
	}
	return nil
}

func SubmitWithRetry(txBuilder *build.TransactionBuilder, seed string, logger log.Logger) error {
	client := horizon.Client{
		URL:  "https://horizon-testnet.stellar.org",
		HTTP: &http.Client{Timeout: transferTimeout},
	}

	txEnv := txBuilder.Sign(seed)
	txEnvB64, err := txEnv.Base64()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}

	for i := 0; i < retryFailedTxAmount; i++ {
		_, err = client.SubmitTransaction(txEnvB64)
		if err == nil {
			break
		}
	}

	return err
}
