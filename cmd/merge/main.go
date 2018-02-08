// Merge accounts and their balance into a single account.
//
// See Keypairs struct for expected input format.
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

const maxMergeOps = 19

var (
	horizonDomainFlag   = flag.String("address", "https://horizon-testnet.stellar.org", "horizon address")
	publicNetworkFlag   = flag.Bool("pubnet", false, "use public network")
	destinationSeedFlag = flag.String("dest", "", "destination account seed")
	accountsFile        = flag.String("input", "accounts.json", "keypairs input file")
)

type Keypair struct {
	Seed    string `json:"seed"`
	Address string `json:"address"`
}

type Keypairs struct {
	Keypairs []Keypair `json:"keypairs"`
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

	// Get destination account
	destKP := keypair.MustParse(*destinationSeedFlag)

	// Send transactions in batches of 20 signers (Stellar limitation)
	for batchIndex := 0; batchIndex <= (len(keypairs.Keypairs)+1)/(maxMergeOps+1); batchIndex++ {
		batch := keypairs.Keypairs[batchIndex*maxMergeOps : int(math.Min(float64((batchIndex+1)*maxMergeOps), float64(len(keypairs.Keypairs))))]

		signers := []string{destKP.(*keypair.Full).Seed()}

		// Generate account merge operations
		var ops []build.TransactionMutator
		for i, kp := range batch {
			level.Debug(logger).Log("account_index", batchIndex*maxMergeOps+i, "msg", "adding merge operation")

			ops = append(
				ops,
				build.AccountMerge(
					build.SourceAccount{AddressOrSeed: kp.Address},
					build.Destination{AddressOrSeed: destKP.Address()},
				),
			)

			signers = append(signers, kp.Seed)
		}

		// Set network where the transaction will be submitted to
		var network build.Network
		if *publicNetworkFlag == true {
			network = build.PublicNetwork
		} else {
			network = build.TestNetwork
		}

		// Add transaction submitter source account and network information
		client := horizon.Client{
			URL:  *horizonDomainFlag,
			HTTP: &http.Client{Timeout: 20 * time.Second},
		}
		ops = append(
			[]build.TransactionMutator{
				build.SourceAccount{AddressOrSeed: destKP.(*keypair.Full).Seed()},
				network,
				build.AutoSequence{SequenceProvider: &client},
			},
			ops...,
		)

		// Generate and submit transaction
		txBuilder, err := build.Transaction(ops...)
		if err != nil {
			level.Error(logger).Log("msg", err)
			os.Exit(1)
		}
		txEnv, err := txBuilder.Sign(signers...)
		if err != nil {
			level.Error(logger).Log("msg", err)
			os.Exit(1)
		}
		var txEnvB64 string
		txEnvB64, err = txEnv.Base64()
		if err != nil {
			level.Error(logger).Log("msg", err)
			os.Exit(1)
		}

		level.Info(logger).Log("msg", "submitting transaction")

		_, err = client.SubmitTransaction(txEnvB64)
		if err == nil {
			level.Info(logger).Log("msg", "submit success")
			continue
		}

		GetTxErrorResultCodes(err, logger)
	}
}
