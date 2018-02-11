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
	client  *horizon.Client
	network build.Network

	sourceSeed,
	sourceAddress,
	destinationAddress,

	transferAmount string

	sequenceProvider *sequence.Provider

	// Submitter will close this channel once it finishes submitting transactions.
	// See StartSubmission() function for more information.
	Stopped chan struct{}
}

// New returns a new Submitter.
func New(
	client *horizon.Client,
	network build.Network,
	provider *sequence.Provider,
	source *keypair.Full,
	destination keypair.KP,
	transferAmount string) (*Submitter, error) {

	s := Submitter{
		client:  client,
		network: network,

		sourceSeed:         source.Seed(),
		sourceAddress:      source.Address(),
		destinationAddress: destination.Address(),

		transferAmount: transferAmount,

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
func (s *Submitter) StartSubmission(ctx context.Context, limiter *rate.Limiter, logger log.Logger) {
	logger = log.With(logger, "source_address", s.sourceAddress)

	go func() {
		level.Debug(logger).Log("msg", "starting")

		defer func() {
			close(s.Stopped)
		}()

		for {
			if err := limiter.Wait(ctx); err != nil {
				// Stop submitting if context is canceled,
				// meaning submission should stop.
				if ctx.Err() != nil {
					break
				}
				continue
			}

			s.submit(logger)
		}
	}()
}

// submit submits a transaction to the Stellar network.
// The transaction has the same property on every call:
// Same source and and desitnation addresses, and same amount.
// The only property that changes is the sequence number.
func (s *Submitter) submit(logger log.Logger) error {
	level.Debug(logger).Log("msg", "building transaction")

	txBuilder, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: s.sourceAddress},
		s.network,
		build.AutoSequence{SequenceProvider: s.sequenceProvider},

		build.Payment(
			build.Destination{AddressOrSeed: s.destinationAddress},
			build.NativeAmount{Amount: s.transferAmount},
		),
	)
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
	_, err = s.client.SubmitTransaction(txEnvB64)
	duration := time.Since(start)
	logger = log.With(logger, "response_time", duration, "response_time_nanoseconds", duration.Nanoseconds())

	// Return success if submission was successful
	if err == nil {
		level.Info(logger).Log("status", "success")
		return nil
	}

	// Else log and return error
	errors.GetTxErrorResultCodes(err, log.With(logger, "transaction_status", "failure"))
	return err
}
