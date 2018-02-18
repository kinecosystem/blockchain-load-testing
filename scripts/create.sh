#!/usr/bin/env bash
# helper script for create test accounts

set -x
set -e

HORIZON="${HORIZON:-http://localhost:8000}"
ACCOUNTS="${ACCOUNTS:-600}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"
FUND_AMOUNT="${FUND_AMOUNT:-2}"

PUBNET="${PUBNET:-}"
# PUBNET="${PUBNET:--pubnet}"

FUNDER_SEED="${FUNDER_SEED}"

make build
./create \
    -address $HORIZON \
    -funder $FUNDER_SEED \
    -accounts $ACCOUNTS \
    -amount $FUND_AMOUNT \
    -output $ACCOUNTS_FILE \
    $PUBNET
