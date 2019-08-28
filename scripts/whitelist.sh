#!/usr/bin/env bash
# helper script for adding addresses to the whitelist

set -x
set -e


make build

./whitelist \
    -horizon $HORIZON \
    -passphrase $NETWORK_PASSPHRASE \
    -whitelist-seed $WHITELIST_SEED \
    -address-to-whitelist $ADDRESS_TO_WHITELIST
