package main

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

// LogBalance logs given account's balances.
func LogBalance(account *horizon.Account, logger log.Logger) {
	for _, balance := range account.Balances {
		level.Info(logger).Log("balance", balance.Balance, "asset_type", balance.Asset.Type)
	}
}

// LogBalances logs given accounts' balances.
func LogBalances(client horizon.ClientInterface, keypairs []keypair.KP, logger log.Logger) {
	var wg sync.WaitGroup
	for _, kp := range keypairs {
		wg.Add(1)
		go func(kp keypair.KP, logger log.Logger) {
			defer wg.Done()

			l := log.With(logger, "address", kp.Address())

			if kp != nil {
				acc, err := client.LoadAccount(kp.Address())
				if err != nil {
					level.Error(l).Log("msg", err)
				}
				LogBalance(&acc, l)
			}
		}(kp, logger)
	}
	wg.Wait()
}
