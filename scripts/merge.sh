#!/usr/bin/env bash
# helper script for merging test accounts

set -x
set -e

HORIZON="${HORIZON:-http://localhost:8000}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"

PUBNET="${PUBNET:-}"
# PUBNET="${PUBNET:--pubnet}"

DEST_SEED="${DEST_SEED}"

make build
./merge \
    -address $HORIZON \
    -input $ACCOUNTS_FILE \
    -dest $DEST_SEED \
    $PUBNET
