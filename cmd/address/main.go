// Get account address from given seed.
package main

import (
	"flag"
	"fmt"

	"github.com/kinecosystem/go/keypair"
)

var seedFlag = flag.String("seed", "", "seed")

func main() {
	flag.Parse()

	kp := keypair.MustParse(*seedFlag)
	fmt.Println(kp.Address())
}
