#!/usr/bin/env bash
# helper script for setting arguments

set -x
set -e

RATE=$1
if [ "$RATE" == "" ]; then
	RATE=0
fi

NATIVE=true
if [ "$2" != "" ]; then
	NATIVE=$2
fi

HORIZON1=$HORIZON
HORIZON2=$HORIZON1
HORIZON3=$HORIZON1
HORIZON4=$HORIZON1
HORIZON5=$HORIZON1
HORIZON6=$HORIZON1
HORIZON7=$HORIZON1

DEBUG="${DEBUG:-true}"
PASSPHRASE="${PASSPHRASE:-"Kin Scaling ; March 2019"}"
LOG="${LOG:-loadtest.log}"
ACCOUNTS_FILE="${ACCOUNTS_FILE:-accounts.json}"
TX_AMOUNT="${TX_AMOUNT:-0.0001}"
OPS_PER_TX="${OPS_PER_TX:-1}"
TIME_LENGTH="${TIME_LENGTH:-20}"
RATE="${RATE:-0}"  # zero disables rate limiting
BURST="${BURST:-100}"

DEST_ACCOUNT="${DEST_ACCOUNT:-dest.json}"


./loadtest \
    -native=$NATIVE \
    -debug=$DEBUG \
    -address $HORIZON1 \
    -address $HORIZON2 \
    -address $HORIZON3 \
    -address $HORIZON4 \
    -address $HORIZON5 \
    -address $HORIZON6 \
    -address $HORIZON7 \
    -passphrase "$PASSPHRASE" \
    -log $LOG \
    -accounts $ACCOUNTS_FILE \
    -txamount $TX_AMOUNT \
    -ops $OPS_PER_TX \
    -length $TIME_LENGTH \
    -submitters $SUBMITTERS \
    -rate $RATE \
    -burst $BURST \
    -dest $DEST_ACCOUNT \
    -whitelisted-account-seed $WHITELISTED_SEED
