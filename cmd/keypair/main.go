// Generate new keypair (Seed, Address) and log to stdout.
package main

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/kinecosystem/go/keypair"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "time", log.DefaultTimestampUTC())

	kp, err := keypair.Random()
	if err != nil {
		logger.Log("msg", err)
		return
	}

	logger.Log("msg", "keypair created", "address", kp.Address(), "seed", kp.Seed())
}
