// Package submitter implements a single transaction submitter to the Stellar network.
// Multiple Submitter instances are used concurrently for load testing the network.
package submitter

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"golang.org/x/time/rate"

	"github.com/kinfoundation/stellar-load-testing/cmd/loadtest/errors"
	"github.com/kinfoundation/stellar-load-testing/cmd/loadtest/sequence"
)

// Submitter continuously submits transactions to the Stellar network according to a rate limiter.
//
// The transactions always consist of a single payment operation to a predefined destination address.
type Submitter struct {
	clients []horizon.Client
	network build.Network

	sourceSeed,
	sourceAddress,

	transferAmount string

	destinationAddresses []keypair.KP

	// Amount of payment operations per transaction.
	opsPerTx int

	sequenceProvider *sequence.Provider

	// Submitter will close this channel once it finishes submitting transactions.
	// See StartSubmission() function for more information.
	Stopped chan struct{}
}

// New returns a new Submitter.
func New(
	clients []horizon.Client,
	network build.Network,
	provider *sequence.Provider,
	source *keypair.Full,
	destination []keypair.KP,
	transferAmount string,
	opsPerTx int) (*Submitter, error) {

	s := Submitter{
		clients: clients,
		network: network,

		sourceSeed:           source.Seed(),
		sourceAddress:        source.Address(),
		destinationAddresses: destination,

		transferAmount: transferAmount,

		opsPerTx: opsPerTx,

		sequenceProvider: provider,

		Stopped: make(chan struct{}),
	}

	// Load and cache sequence number for given source account.
	_, err := s.sequenceProvider.SequenceForAccount(s.sourceAddress)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// StartSubmission continously submits transactions to the network using the given rate limiter.
func (s *Submitter) StartSubmission(ctx context.Context, limiter *rate.Limiter, logger log.Logger, native bool) {
	logger = log.With(logger, "source_address", s.sourceAddress)

	go func() {
		level.Debug(logger).Log("msg", "starting")

		defer func() {
			close(s.Stopped)
		}()

		if native {
			level.Debug(logger).Log("msg", "using native asset")
		} else {
			level.Debug(logger).Log("msg", "using non-native asset")
		}

		destIndex := 0
		clientIndex := 0

		for {
			destIndex++
			if destIndex == len(s.destinationAddresses) {
				destIndex = 0
			}

			clientIndex++
			if clientIndex == len(s.clients) {
				clientIndex = 0
			}

			if err := limiter.Wait(ctx); err != nil {
				// Stop submitting if context is canceled,
				// meaning submission should stop.
				if ctx.Err() != nil {
					break
				}
				continue
			}

			s.submit(logger, destIndex, native, clientIndex)
		}
	}()
}

// submit submits a transaction to the Stellar network.
// The transaction has the same property on every call:
// Same source and and desitnation addresses, and same amount.
// The only property that changes is the sequence number.
func (s *Submitter) submit(logger log.Logger, destIndex int, native bool, clientIndex int) error {
	level.Debug(logger).Log("msg", "building transaction", "ops_per_tx", s.opsPerTx)

	ops := append(
		[]build.TransactionMutator{},

		build.SourceAccount{AddressOrSeed: s.sourceAddress},
		s.network,
		build.AutoSequence{SequenceProvider: s.sequenceProvider},
	)

	for i := 0; i < s.opsPerTx; i++ {
		var amount build.PaymentMutator

		if native {
			amount = build.NativeAmount{Amount: s.transferAmount}
		} else {
			amount = build.CreditAmount{"KIN", "GBSJ7KFU2NXACVHVN2VWQIXIV5FWH6A7OIDDTEUYTCJYGY3FJMYIDTU7", s.transferAmount}
		}

		ops = append(ops, build.Payment(build.Destination{AddressOrSeed: s.destinationAddresses[destIndex].Address()}, amount))
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

	txEnv, err := txBuilder.Sign(s.sourceSeed)
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}

	var txEnvB64 string
	txEnvB64, err = txEnv.Base64()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return err
	}

	level.Info(logger).Log("msg", "submitting transaction")

	start := time.Now()
	_, err = s.clients[clientIndex].SubmitTransaction(txEnvB64)
	duration := time.Since(start)
	logger = log.With(logger, "response_time", duration, "response_time_nanoseconds", duration.Nanoseconds())

	// Return success if submission was successful
	if err == nil {
		l := log.With(logger, "transaction_status", "success")

		if _, err := s.sequenceProvider.IncrementSequence(s.sourceAddress); err != nil {
			level.Error(l).Log("sequence_provider_error", err)
			return nil
		}

		level.Info(l).Log()

		return nil
	}

	// Logs errors and set the current sequence number from Horizon instead of local cache
	// if transaction failed due to bad sequence number.
	code := errors.GetTxErrorResultCodes(err, log.With(logger, "transaction_status", "failure"))
	if code != nil && code.TransactionCode == "tx_bad_seq" {
		if _, err := s.sequenceProvider.LoadSequenceWithClient(s.sourceAddress); err != nil {
			level.Error(logger).Log("msg", err)
		}
	}
	return err
}
