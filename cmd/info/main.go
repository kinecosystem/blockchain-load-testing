// Print account info for given address.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/stellar/go/clients/horizon"
)

var (
	horizonDomainFlag = flag.String("horizon", "https://horizon-testnet.stellar.org", "horizon address")
	addressFlag       = flag.String("address", "", "account address")
)

func main() {
	flag.Parse()

	client := horizon.Client{
		URL:  *horizonDomainFlag,
		HTTP: &http.Client{Timeout: 5 * time.Second},
	}

	account, err := client.LoadAccount(*addressFlag)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(&account)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
