#!/usr/bin/env bash
# helper script for merging test accounts

set -x
set -e

HORIZON="${HORIZON:-https://horizon-testnet.stellar.org}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"

DEST_SEED="SAPJIKONFL75FU4NVPW7GWXFWBIM2N4DRSIAV2HLFOG2UGO26XVZGNA2"

go run cmd/merge/*.go \
    -address $HORIZON \
    -input $ACCOUNTS_FILE \
    -dest $DEST_SEED
