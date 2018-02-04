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

// Provider provider sequence numbers for Stellar transactions,
// with caching. This saves on executing multiple requests to an Horizon
// instance for fetching an account's sequence number.
//
// Note this package assumes you are using no more than a single provider for
// an account. Otherwise the returned sequence number will be incorrect
// and out of sync.
type Provider struct {
	build.SequenceProvider

	locker sync.Locker
	client horizon.ClientInterface

	// Account address to sequence number mapping.
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

// SequenceForAccount returns the next available sequence number for given account
// address by fetching the current sequence and increasing it by one.
func (p *Provider) SequenceForAccount(address string) (xdr.SequenceNumber, error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	seq, err := p.getAndCache(address)
	if err != nil {
		return 0, err
	}

	return seq, nil
	// newSeq := seq + 1
	// p.sequences[address] = newSeq
	// return newSeq, nil
}

// GetAndCache fetches, caches and returns the current sequence number for the given
// account address.
//
// Note this function can act as a "cache warmup" for an account,
// since it loads the sequence number from Horizon when called for the first
// time for a particular account.
func (p *Provider) GetAndCache(address string) (xdr.SequenceNumber, error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	return p.getAndCache(address)
}

func (p *Provider) getAndCache(address string) (xdr.SequenceNumber, error) {
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
