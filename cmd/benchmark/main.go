package main

import (
	"os"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"

	"github.com/kinfoundation/stellar-benchmark/src/account"
	"github.com/kinfoundation/stellar-benchmark/src/transaction"
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

func logBalance(account *horizon.Account, logger log.Logger) {
	for _, balance := range account.Balances {
		level.Info(logger).Log("balance", balance.Balance, "asset_type", balance.Asset.Type)
	}
}

func logBalances(keypairs []keypair.KP, logger log.Logger) {
	for i, kp := range keypairs {
		l := log.With(logger, "account_index", i)
		if kp != nil {
			acc := account.Get(kp.Address(), l)
			logBalance(acc, l)
		}
	}
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())
	// logger = log.With(logger, "time", log.DefaultTimestampUTC())
	// logger = log.With(logger, "caller", log.Caller(3))

	funderKP := keypair.MustParse(funderSeed)
	funderAccount := account.Get(funderKP.Address(), logger)
	logBalance(funderAccount, log.With(logger, "msg", "funder account info", "address", funderKP.Address()[:5], "seed", funderSeed))

	keypairs, err := account.Create(funderKP.(*keypair.Full), accountsNum, fundAmount, logger)
	if err != nil {
		return
	}
	logBalances(keypairs, logger)

	var wg sync.WaitGroup
	for i, kp := range keypairs {
		wg.Add(1)
		go func(kp keypair.KP) {
			defer wg.Done()

			l := log.With(logger, "account_index", i)

			for j, other := range keypairs {
				if kp.Address() == other.Address() {
					continue
				}

				l = log.With(l, "iteration", j)

				if err := transaction.Transfer(kp, other, transferAmount, l); err != nil {
					level.Error(logger).Log("msg", err)
					os.Exit(1)
				}
			}
		}(kp)
	}
	wg.Wait()

	logBalances(keypairs, logger)
}
