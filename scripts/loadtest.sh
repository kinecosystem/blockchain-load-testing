#!/usr/bin/env bash
# helper script for setting arguments

set -x
set -e

SUBMITTERS=$1
if [ "$SUBMITTERS" == "" ]; then
	SUBMITTERS=0
fi

NATIVE=true
if [ "$2" != "" ]; then
	NATIVE=$2
fi

DEBUG="${DEBUG:-true}"
HORIZON="${HORIZON:-http://localhost:8000}"
PASSPHRASE="${PASSPHRASE:-"private testnet"}"
LOG="${LOG:-loadtest.log}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"
TX_AMOUNT="${TX_AMOUNT:-0.0001}"
OPS_PER_TX="${OPS_PER_TX:-1}"
TIME_LENGTH="${TIME_LENGTH:-20}"
RATE="${RATE:-0}"  # zero disables rate limiting
BURST="${BURST:-100}"

DEST_ACCOUNT="${DEST_ACCOUNT:-dest.json}"

make build
./loadtest \
    -native=$NATIVE \
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
