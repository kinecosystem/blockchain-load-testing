#!/usr/bin/env bash
# helper script for setting arguments

set -x
set -e

DEBUG="${DEBUG:-true}"
HORIZON="${HORIZON:-https://horizon-testnet.stellar.org}"
LOG="${LOG:-loadtest.log}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"
TX_AMOUNT="${TX_AMOUNT:-0.0001}"
TIME_LENGTH="${TIME_LENGTH:-10800}"  # three hours
SUBMITTERS="${SUBMITTERS:-120}"
RATE="${RATE:-2}"
BURST="${BURST:-2}"

PUBNET="${PUBNET:-}"
# PUBNET="${PUBNET:--pubnet}"

DEST_ACCOUNT="${DEST_ACCOUNT}"

./loadtest \
    -debug=$DEBUG \
    -address $HORIZON \
    -log $LOG \
    -accounts $ACCOUNTS_FILE \
    -txamount $TX_AMOUNT \
    -length $TIME_LENGTH \
    -submitters $SUBMITTERS \
    -rate $RATE \
    -burst $BURST \
    -dest $DEST_ACCOUNT \
    $PUBNET
