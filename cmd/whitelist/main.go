// Add addresses to whitelist.
package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kinecosystem/go/build"
	"github.com/kinecosystem/go/clients/horizon"
	"github.com/kinecosystem/go/keypair"
)

const ClientTimeout = 30 * time.Second

// Support repeating string flags
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join((*i)[:], ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	horizonAddressFlag      = flag.String("horizon", "", "horizon address")
	networkPassphraseFlag   = flag.String("passphrase", "", "network passphrase")
	whitelistSeedFlag       = flag.String("whitelist-seed", "", "whitelist seed; address is dervied from it")
	addressToWhitelistFlags arrayFlags
)

func whitelist(client *horizon.Client, network *build.Network, whitelistKP *keypair.Full, kpsToWhitelist []keypair.KP, logger log.Logger) error {
	level.Debug(logger).Log("msg", "building transaction")

	ops := append(
		[]build.TransactionMutator{},

		build.SourceAccount{AddressOrSeed: whitelistKP.Seed()},
		network,
		build.AutoSequence{SequenceProvider: client},
	)

	for _, kp := range kpsToWhitelist {
		hint := kp.Hint()
		ops = append(ops, build.SetData(kp.Address(), hint[:]))
	}

	txBuilder, err := build.Transaction(ops...)
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}

	txHash, err := txBuilder.HashHex()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}
	logger = log.With(logger, "tx_hash", txHash)

	txEnv, err := txBuilder.Sign(whitelistKP.Seed())

	var txEnvB64 string
	txEnvB64, err = txEnv.Base64()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}

	level.Info(logger).Log(
		"msg", "submitting transaction",
		"tx_env_b64", txEnvB64,
	)

	start := time.Now()
	_, err = client.SubmitTransaction(txEnvB64)
	duration := time.Since(start)
	logger = log.With(logger, "response_time", duration)

	// Return success if submission was successful
	if err == nil {
		l := log.With(logger, "transaction_status", "success")
		level.Info(l).Log()
		return nil
	}

	level.Error(logger).Log("msg", err)

	// Logs errors and set the current sequence number from Horizon instead of local cache
	// if transaction failed due to bad sequence number.
	getTxErrorResultCodes(err, log.With(logger, "transaction_status", "failure"))

	return err
}

func getTxErrorResultCodes(err error, logger log.Logger) *horizon.TransactionResultCodes {
	level.Error(logger).Log("msg", err)
	switch e := err.(type) {
	case *horizon.Error:
		code, err := e.ResultCodes()
		if err != nil {
			level.Error(logger).Log("msg", "failed to extract result codes from horizon response")
			return nil
		}
		level.Error(logger).Log("code", code.TransactionCode)
		for i, opCode := range code.OperationCodes {
			level.Error(logger).Log("opcode_index", i, "opcode", opCode)
		}

		return code
	}
	return nil
}

func init() {
	flag.Var(&addressToWhitelistFlags, "address-to-whitelist", "Address to add to whitelist; Flag can be repeated multiple times for whitelisting multiple addresses")
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
		URL:  *horizonAddressFlag,
		HTTP: &http.Client{Timeout: ClientTimeout},
	}

	network := build.Network{*networkPassphraseFlag}

	// Create whitelist Keypair struct
	whitelistKP := keypair.MustParse(*whitelistSeedFlag)
	level.Info(logger).Log("msg", "whitelist account info", "address", whitelistKP.Address()[:5], "seed", *whitelistSeedFlag)

	// Create address to whitelist Keypair slice
	var kpsToWhitelist []keypair.KP
	for _, address := range addressToWhitelistFlags {
		kpsToWhitelist = append(kpsToWhitelist, keypair.MustParse(address))
	}

	whitelist(&client, &network, whitelistKP.(*keypair.Full), kpsToWhitelist, logger)
}
