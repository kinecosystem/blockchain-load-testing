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

pipenv sync; pipenv run python create.py --channel-seeds-file channel_seed --accounts $AMOUNT_OF_SOURCE_ACCOUNTS --passphrase "$PASS" --horizon $HORIZON --source-account $FUNDER_SEED --json-output True >../../accounts.json
pipenv run python create.py --channel-seeds-file channel_seed --accounts $AMOUNT_OF_DESTINATION_ACCOUNTS --passphrase "$PASS" --horizon $HORIZON --source-account $FUNDER_SEED --json-output True >../../dest.json
