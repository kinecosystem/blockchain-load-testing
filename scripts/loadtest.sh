#!/usr/bin/env bash
# helper script for setting arguments

set -x
set -e

COUNT=$1

if [ "$COUNT" == "" ]; then
	echo "Usage: $0 <COUNT>"
	exit 1
fi

SUBMITTERS=$2
if [ "$SUBMITTERS" == "" ]; then
	SUBMITTERS="${COUNT}"
fi

DEBUG="${DEBUG:-true}"
HORIZON="${HORIZON:-http://localhost:8000}"
PASSPHRASE="${PASSPHRASE:-"private testnet"}"
LOG="${LOG:-loadtest.log}"
ACCOUNTS_FILE="${COUNT}.json"
TX_AMOUNT="${TX_AMOUNT:-0.0001}"
OPS_PER_TX="${OPS_PER_TX:-1}"
TIME_LENGTH="${TIME_LENGTH:-20}"
RATE="${RATE:-0}"  # zero disables rate limiting
BURST="${BURST:-100}"

DEST_ACCOUNT="${DEST_ACCOUNT}"

make build
./loadtest \
    -debug=$DEBUG \
    -address $HORIZON \
    -passphrase "$PASSPHRASE" \
    -log $LOG \
    -accounts $ACCOUNTS_FILE \
    -txamount $TX_AMOUNT \
    -ops $OPS_PER_TX \
    -length $TIME_LENGTH \
    -submitters $SUBMITTERS \
    -rate $RATE \
    -burst $BURST \
    -dest $DEST_ACCOUNT
