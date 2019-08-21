#!/usr/bin/env bash
# bash script helper to create accounts (using ../cmd/python)
# Parameters taken from  vars.sh
set -x
set -e
. ../../vars.sh

#CHANNEL_SEED=$1
#AMOUNT_OF_SOURCE_ACCOUNTS=
#AMOUNT_OF_DESTINATION_ACCOUNTS=
#PASSPHRASE=$3
#HORIZON=$4
#FUNDER_SEED=$5
echo $CHANNEL_SEED > channel_seed
PASS=$PASSPHRASE

# create whitelist account using root account
curl \
    -X POST \
    -H 'multipart/form-data' \
    -d 'tx=AAAAABTvssx9E+cge57y3TtzZCW5vuX/zhB3f5DRU7uNU6P6AAAAZAAAAAAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAJYQgQN++M4kAnwTTDZPyKCvoievyTVX1O7S/GKIRUdUAAAAAAABhqAAAAAAAAAAAA==' \
    'http://ip-horizon-test-1.test.kin/transactions'

pipenv sync; pipenv run python create.py --channel-seeds-file channel_seed --accounts $AMOUNT_OF_SOURCE_ACCOUNTS --passphrase "$PASS" --horizon $HORIZON --source-account $FUNDER_SEED --json-output True >../../accounts.json

pipenv run python create.py --channel-seeds-file channel_seed --accounts $AMOUNT_OF_DESTINATION_ACCOUNTS --passphrase "$PASS" --horizon $HORIZON --source-account $FUNDER_SEED --json-output True >../../dest.json

# whitelist one address, used to whitelist txs in load test
../../resources/whitelist \
    -horizon $HORIZON \
    -passphrase $PASS \
    -whitelist-seed $WHITELIST_SEED \
    -address-to-whitelist $WHITELISTED_SEED
