// Read accounts JSON file and print account balances to screen.
//
// See Keypairs struct for expected input format.
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
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

var (
	horizonDomainFlag = flag.String("address", "https://horizon-testnet.stellar.org", "horizon address")
	accountsFile      = flag.String("input", "accounts.json", "keypairs input file")
)

type Keypair struct {
	Seed string `json:"seed"`
}

type Keypairs struct {
	Keypairs []Keypair `json:"keypairs"`
}

func logBalance(account *horizon.Account, logger log.Logger) {
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

	// Read accounts file
	b, err := ioutil.ReadFile(*accountsFile)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
	var keypairs Keypairs
	err = json.Unmarshal(b, &keypairs)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	// Log accounts
	client := horizon.Client{
		URL:  *horizonDomainFlag,
		HTTP: &http.Client{Timeout: 5 * time.Second},
	}
	for _, kpObj := range keypairs.Keypairs {
		kp := keypair.MustParse(kpObj.Seed)
		acc, err := client.LoadAccount(kp.Address())
		if err != nil {
			level.Error(logger).Log("msg", err)
			os.Exit(1)
		}
		for _, balance := range acc.Balances {
			level.Info(logger).Log("address", kp.Address()[:5], "balance", balance.Balance, "asset_type", balance.Asset.Type)
		}
	}
}
