// Found an account (on testnet) using friendbot.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kinecosystem/go/clients/horizon"
	"github.com/kinecosystem/go/keypair"
)

var (
	horizonDomainFlag = flag.String("address", "https://horizon-testnet.stellar.org", "horizon address")
)

func Fund(horizonAddr string, kp keypair.KP, logger log.Logger) error {
	l := log.With(logger, "address", kp.Address()[:5])

	level.Info(l).Log("msg", "sending funding request", "address", kp.Address()[:5])

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Get(fmt.Sprintf("%s/friendbot?addr=%s", horizonAddr, kp.Address()))
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			level.Error(l).Log("msg", err)
		}
	}()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	level.Debug(l).Log("msg", string(data))

	if res.StatusCode == http.StatusBadRequest {
		return errors.New("funding failure")
	}

	level.Info(l).Log("msg", "funding success")

	return nil
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	// logger = level.NewFilter(logger, level.AllowDebug())
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "time", log.DefaultTimestampUTC())
	// logger = log.With(logger, "caller", log.Caller(3))

	flag.Parse()

	kp, err := keypair.Random()
	if err != nil {
		level.Error(logger).Log("msg", err)
		return
	}

	level.Info(logger).Log("msg", "keypair created", "address", kp.Address(), "seed", kp.Seed())

	if err := Fund(*horizonDomainFlag, kp, logger); err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	l := log.With(logger, "address", kp.Address()[:5])

	client := horizon.Client{
		URL:  *horizonDomainFlag,
		HTTP: &http.Client{Timeout: 10 * time.Second},
	}

	account, err := client.LoadAccount(kp.Address())
	if err != nil {
		level.Error(l).Log("msg", err)
		os.Exit(1)
	}

	for _, balance := range account.Balances {
		level.Info(logger).Log("balance", balance.Balance, "asset_type", balance.Asset.Type)
	}
}
