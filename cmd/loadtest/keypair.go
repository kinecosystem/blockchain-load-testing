package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/stellar/go/keypair"
)

type KeypairJSON struct {
	Seed    string `json:"seed"`
	Address string `json:"address"`
}

type KeypairsJSON struct {
	Keypairs []KeypairJSON `json:"keypairs"`
}

// InitKeypairs is a helper function that reads worker accounts file
// and returns their keypair.KP objects.
func InitKeypairs(path string) ([]keypair.KP, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var keypairsJSON KeypairsJSON
	err = json.Unmarshal(b, &keypairsJSON)
	if err != nil {
		return nil, err
	}

	keypairs := make([]keypair.KP, len(keypairsJSON.Keypairs))
	for i := 0; i < len(keypairsJSON.Keypairs); i++ {
		keypairs[i] = keypair.MustParse(keypairsJSON.Keypairs[i].Seed)
	}

	return keypairs, nil
}
