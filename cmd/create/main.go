// Create and fund multiple accounts.
//
// Seed and address keypairs are dumped to JSON file.
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

const (
	loadAccountTimeout = 5 * time.Second
)

var (
	horizonDomainFlag = flag.String("address", "https://horizon-testnet.stellar.org", "horizon address")
	publicNetworkFlag = flag.Bool("pubnet", false, "use public network")
	funderSeedFlag    = flag.String("funder", "", "funder seed")
	accountsNumFlag   = flag.Int("accounts", 0, "amount of accounts to create and fund")
	fundAmountFlag    = flag.String("amount", "0.5", "funding amount for each account")
	keypairsFile      = flag.String("output", "accounts.json", "keypairs output file")
)

type Keypair struct {
	Seed    string `json:"seed"`
	Address string `json:"address"`
}

type Keypairs struct {
	Keypairs []Keypair `json:"keypairs"`
}

func logBalance(account *horizon.Account, logger log.Logger) {
	for _, balance := range account.Balances {
		level.Info(logger).Log("balance", balance.Balance, "asset_type", balance.Asset.Type)
	}
}

func logBalances(client horizon.ClientInterface, keypairs []keypair.KP, logger log.Logger) {
	for i, kp := range keypairs {
		l := log.With(logger, "account_index", i)

		if kp != nil {
			acc, err := client.LoadAccount(kp.Address())
			if err != nil {
				level.Error(l).Log("msg", err)
				continue
			}
			logBalance(&acc, l)
		}
	}
}

func main() {
	flag.Parse()

	// Initialize logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())

	// Log execution time
	start := time.Now()
	defer func(logger log.Logger) {
		level.Info(logger).Log("execution_time", time.Since(start))
	}(logger)

	client := horizon.Client{
		URL:  *horizonDomainFlag,
		HTTP: &http.Client{Timeout: loadAccountTimeout},
	}

	// Get funding account
	funderKP := keypair.MustParse(*funderSeedFlag)
	funderAccount, err := client.LoadAccount(funderKP.Address())
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
	logBalance(&funderAccount, log.With(logger, "msg", "funder account info", "address", funderKP.Address()[:5], "seed", *funderSeedFlag))

	var network build.Network
	if *publicNetworkFlag == true {
		network = build.PublicNetwork
	} else {
		network = build.TestNetwork
	}

	// Create and fund accounts
	keypairs, err := Create(*horizonDomainFlag, network, funderKP.(*keypair.Full), *accountsNumFlag, *fundAmountFlag, logger)
	if err != nil {
		os.Exit(1)
	}
	logBalances(&client, keypairs, logger)

	// Write the seeds of the created accounts to file.
	keypairsOut := Keypairs{Keypairs: make([]Keypair, 0)}
	for _, kp := range keypairs {
		keypairsOut.Keypairs = append(
			keypairsOut.Keypairs,
			Keypair{
				Seed:    kp.(*keypair.Full).Seed(),
				Address: kp.Address(),
			},
		)
	}

	keypairsData, err := json.Marshal(keypairsOut)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(*keypairsFile, keypairsData, 0644); err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}
