package main

import (
	"flag"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"

	"github.com/kinfoundation/stellar-benchmark/src/account"
	"github.com/kinfoundation/stellar-benchmark/src/transaction"
)

const (
	funderSeed = "SCY4OBPZQCCSJLC22VD47JXCLJKXQKM5VR56RY2CPY4TPCOK37G57QJA"

	accountsNum = 100

	fundAmount     = "30"
	transferAmount = "0.01"
)

var (
	horizonDomainFlag = flag.String("address", "https://horizon-testnet.stellar.org", "horizon address")
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
			acc, err := account.Get(*horizonDomainFlag, kp.Address(), l)
			if err != nil {
				os.Exit(1)
			}
			logBalance(acc, l)
		}
	}
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())
	// logger = log.With(logger, "time", log.DefaultTimestampUTC())
	// logger = log.With(logger, "caller", log.Caller(3))

	flag.Parse()

	level.Info(logger).Log("horizon_address", *horizonDomainFlag)

	funderKP := keypair.MustParse(funderSeed)
	funderAccount, err := account.Get(*horizonDomainFlag, funderKP.Address(), logger)
	if err != nil {
		os.Exit(1)
	}
	logBalance(funderAccount, log.With(logger, "msg", "funder account info", "address", funderKP.Address()[:5], "seed", funderSeed))

	keypairs, err := account.Create(*horizonDomainFlag, funderKP.(*keypair.Full), accountsNum, fundAmount, logger)
	if err != nil {
		os.Exit(1)
	}
	logBalances(keypairs, logger)

	var wg sync.WaitGroup
	for i, kp := range keypairs {
		wg.Add(1)
		go func(i int, kp keypair.KP) {
			defer wg.Done()

			l := log.With(logger, "account_index", i)

			for j, other := range keypairs {
				if kp.Address() == other.Address() {
					continue
				}

				err := transaction.Transfer(*horizonDomainFlag, kp, other, transferAmount, log.With(l, "iteration", j))
				if err != nil {
					level.Error(logger).Log("msg", err)
					os.Exit(1)
				}
			}
		}(i, kp)
	}
	wg.Wait()

	logBalances(keypairs, logger)
}
