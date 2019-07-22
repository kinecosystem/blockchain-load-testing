# Stellar Load Testing

Code for load testing the Stellar network, written in Go.

Kin are looking to migrate away from Ethereum and onto a more predictable blockchain in terms of block time and fees.
Stellar is a good candidate, and as part of the process we're testing the network indeed stand up to the performance
it claims it has.

## Build

This application was developed using Go 1.9, though earlier versions may work as well.
Go expects a specific folder formation: /<Local-path>/work/go/src/github.com/kinfoundation  
and environment variable:  
`export GOPATH=/<Local-path>/work/go`

Run from /<Local-path>/work/go/src/github.com/kinfoundation/stellar-load-testing:

```bash
# download glide
make glide

# install dependencies
make vendor

# build binaries
make build
```

## Run

[cmd/loadtest/](cmd/loadtest/) is the main application used for load testing.
Additional helper apps can be found in [cmd/](cmd/).
Check code comments in each for more information.

```bash
# fund test account, testnet only
go run cmd/friendbot

# create and fund test accounts
# see scripts source for default flags
./scripts/create.sh

# run load test
./scripts/loadtest.sh

# merge test accounts back into a single account
./scripts/merge.sh
```

## Generate Reports

The [reports/](reports/) directory contains short Python scripts that parse load test logs
and generate CSV files, ready for charting with a spreadsheet editor like Google Spreadsheets or Excel.
