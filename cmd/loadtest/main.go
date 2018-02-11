// Load test the Stellar network.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"golang.org/x/time/rate"

	"github.com/kinfoundation/stellar-load-testing/cmd/loadtest/sequence"
	"github.com/kinfoundation/stellar-load-testing/cmd/loadtest/submitter"
)

// ClientTimeout is the Horizon HTTP request timeout.
const ClientTimeout = 60 * time.Second

var (
	debugFlag              = flag.Bool("debug", false, "enable debug log level")
	horizonDomainFlag      = flag.String("address", "https://horizon-testnet.stellar.org", "horizon address")
	publicNetworkFlag      = flag.Bool("pubnet", false, "use public network")
	logFileFlag            = flag.String("log", "loadtest.log", "log file path")
	destinationAddressFlag = flag.String("dest", "", "destination account address")
	accountsFileFlag       = flag.String("accounts", "accounts.json", "accounts keypairs input file")
	transactionAmountFlag  = flag.String("txamount", "0.00001", "transaction amount")
	testTimeLengthFlag     = flag.Int("length", 60, "test length in seconds")
	numSubmittersFlag      = flag.Int("submitters", 3, "amount of concurrent submitters")
	txsPerSecondFlag       = flag.Float64("rate", 10, "transaction rate limit in seconds")
	burstLimitFlag         = flag.Int("burst", 3, "burst rate limit")
)

// Run is the main function of this application. It returns a status exit code for main().
func Run() int {
	flag.Parse()

	switch {
	case *destinationAddressFlag == "":
		fmt.Println("-dest flag not set")
		return 1
	}

	// Init logger
	logFile, err := os.OpenFile(*logFileFlag, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := InitLoggers(logFile, *debugFlag)

	// Load destination account address
	destKP, err := keypair.Parse(*destinationAddressFlag)
	if err != nil {
		level.Error(logger).Log("msg", err)
		return 1
	}

	// Load submitter account keypairs
	keypairs, err := InitKeypairs(*accountsFileFlag)
	if err != nil {
		level.Error(logger).Log("msg", err)
		return 1
	}

	client := horizon.Client{
		URL:  *horizonDomainFlag,
		HTTP: &http.Client{Timeout: ClientTimeout},
	}

	LogBalances(&client, keypairs, logger)

	// Init rate limiter
	limiter := rate.NewLimiter(rate.Limit(*txsPerSecondFlag), *burstLimitFlag)

	// Create top-level context. Will be sent to submitter goroutines for stopping them
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel the context if not done so already when test is complete.

	var network build.Network
	if *publicNetworkFlag == true {
		network = build.PublicNetwork
	} else {
		network = build.TestNetwork
	}

	// Generate workers for submitting operations.
	submitters := make([]*submitter.Submitter, *numSubmittersFlag)
	sequenceProvider := sequence.New(&client, logger)
	for i := 0; i < *numSubmittersFlag; i++ {
		level.Debug(logger).Log("msg", "creating submitter", "submitter_index", i)
		submitters[i], err = submitter.New(&client, network, sequenceProvider, keypairs[i].(*keypair.Full), destKP, *transactionAmountFlag)
		if err != nil {
			level.Error(logger).Log("msg", err, "submitter_index", i)
			return 1
		}
	}

	// Start transaction submission
	startTime := time.Now()
	for i := 0; i < *numSubmittersFlag; i++ {
		submitters[i].StartSubmission(ctx, limiter, logger)
	}

	// Listen for OS signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Stop when timer is up or when a signal is caught
	select {
	case <-time.After(time.Duration(*testTimeLengthFlag) * time.Second):
		level.Info(logger).Log("msg", "test time reached")
		break
	case s := <-done:
		level.Info(logger).Log("msg", "received signal", "type", s)
		break
	}
	level.Info(logger).Log("msg", "closing")

	// Stop all submitters
	cancel()
	var wg sync.WaitGroup
	for i, s := range submitters {
		wg.Add(1)
		go func(i int, s *submitter.Submitter) {
			defer wg.Done()
			<-submitters[i].Stopped
		}(i, s)
	}
	wg.Wait()

	level.Info(logger).Log("execution_time", time.Since(startTime))

	// Print destination and test accounts balances
	destAccount, err := client.LoadAccount(destKP.Address())
	LogBalance(&destAccount, logger)

	LogBalances(&client, keypairs, logger)

	return 0
}

func main() {
	os.Exit(Run())
}
