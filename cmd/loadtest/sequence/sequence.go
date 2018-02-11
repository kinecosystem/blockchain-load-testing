// Package sequence implements an in-app sequence provider for the Stellar network,
// independent of Horizon nodes.
package sequence

import (
	"strconv"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
)

// Provider provides sequence numbers for Stellar transactions,
// with local in-app caching. This saves on executing multiple requests to an Horizon
// instance for fetching an account's sequence number.
//
// Note this package assumes you are using no more than a single provider for
// an account. Otherwise the returned sequence number will be incorrect
// and out of sync.
type Provider struct {
	build.SequenceProvider

	locker sync.Locker
	client horizon.ClientInterface

	// Local account sequence number cache
	sequences map[string]xdr.SequenceNumber

	logger log.Logger
}

// New receives an Horizon client and returns a new Provider instance.
func New(c horizon.ClientInterface, logger log.Logger) *Provider {
	return &Provider{
		locker:    &sync.Mutex{},
		client:    c,
		sequences: make(map[string]xdr.SequenceNumber),
		logger:    logger,
	}
}

// SequenceForAccount returns the sequence number for given account using local cache.
func (p *Provider) SequenceForAccount(address string) (xdr.SequenceNumber, error) {
	l := log.With(p.logger, "source_address", address)

	// Fetch sequence number from Horizon if not found in cache.
	seq, ok := p.sequences[address]
	if !ok {
		l = log.With(l, "sequence_provider_source", "horizon client")

		account, err := p.client.LoadAccount(address)
		if err != nil {
			return 0, err
		}

		seqUint, err := strconv.ParseUint(account.Sequence, 10, 64)
		if err != nil {
			return 0, err
		}
		seq = xdr.SequenceNumber(seqUint)

		p.sequences[address] = seq
	} else {
		l = log.With(l, "sequence_provider_source", "local cache")
	}

	level.Debug(l).Log("msg", "sequence number fetched", "sequence_number", seq)

	return seq, nil
}
