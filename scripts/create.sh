#!/usr/bin/env bash
# helper script for create test accounts

set -x
set -e

HORIZON="${HORIZON:-https://horizon-testnet.stellar.org}"
ACCOUNTS="${ACCOUNTS:-600}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"
FUND_AMOUNT="${FUND_AMOUNT:-3}"

FUNDER_SEED="SAPJIKONFL75FU4NVPW7GWXFWBIM2N4DRSIAV2HLFOG2UGO26XVZGNA2"

go run cmd/create/*.go \
    -address $HORIZON \
    -funder $FUNDER_SEED \
    -accounts $ACCOUNTS \
    -amount $FUND_AMOUNT \
    -output $ACCOUNTS_FILE
