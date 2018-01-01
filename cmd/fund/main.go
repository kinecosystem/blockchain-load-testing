package main

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"

	"github.com/kinfoundation/stellar-benchmark/src/account"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	// logger = level.NewFilter(logger, level.AllowDebug())
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "time", log.DefaultTimestampUTC())
	// logger = log.With(logger, "caller", log.Caller(3))

	kp, err := keypair.Random()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return
	}

	level.Info(logger).Log("msg", "keypair created", "address", kp.Address(), "seed", kp.Seed())

	if err := account.Fund(kp, logger); err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	l := log.With(logger, "address", kp.Address()[:5])

	account, err := horizon.DefaultTestNetClient.LoadAccount(kp.Address())
	if err != nil {
		level.Error(l).Log("msg", err)
		os.Exit(1)
	}

	for _, balance := range account.Balances {
		level.Info(logger).Log("balance", balance.Balance, "asset_type", balance.Asset.Type)
	}
}
