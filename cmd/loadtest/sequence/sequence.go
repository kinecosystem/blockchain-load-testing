// Package sequence implements an in-app sequence provider for the Stellar network,
// independent of Horizon nodes.
package sequence

import (
	"strconv"
	"sync"

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
}

// New receives an Horizon client and returns a new Provider instance.
func New(c horizon.ClientInterface) *Provider {
	return &Provider{
		locker:    &sync.Mutex{},
		client:    c,
		sequences: make(map[string]xdr.SequenceNumber),
	}
}

// SequenceForAccount returns the sequence number for given account using local cache.
func (p *Provider) SequenceForAccount(address string) (xdr.SequenceNumber, error) {
	// Fetch sequence number from Horizon if not found in cache.
	seq, ok := p.sequences[address]
	if !ok {
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
	}

	return seq, nil
}
