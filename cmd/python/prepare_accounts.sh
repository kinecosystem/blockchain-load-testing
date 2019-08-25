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

# set base reserve to 0
for core in $(echo $CORE_SERVERS | tr -s " " "\n") ; do
    curl -sf "http://$core:11626/upgrades?basereserve=0&mode=set&upgradetime=1970-01-01T00:00:00Z"
done
# wait for cores to vote and execute upgrade
sleep 15

# set tx set size to 500
for core in $(echo $CORE_SERVERS | tr -s " " "\n") ; do
    # protocol is already at 9, but this fixes a bug which prevents manageData ops
    curl -sf "http://$core:11626/upgrades?protocolversion=9&mode=set&upgradetime=1970-01-01T00:00:00Z"
done
sleep 15

for core in $(echo $CORE_SERVERS | tr -s " " "\n") ; do
    curl -sf "http://$core:11626/upgrades?maxtxsize=500&mode=set&upgradetime=1970-01-01T00:00:00Z"
done

# create whitelist account using funder seed if it doesn't exist yet
if ! $(curl -sf $HORIZON/accounts/$WHITELIST_ADDRESS) ; then
    curl \
        -X POST \
        -H 'multipart/form-data' \
        -d 'tx=AAAAABTvssx9E+cge57y3TtzZCW5vuX/zhB3f5DRU7uNU6P6AAAAZAAAAAAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAJYQgQN++M4kAnwTTDZPyKCvoievyTVX1O7S/GKIRUdUAAAAAAABhqAAAAAAAAAAAA==' \
        "$HORIZON/transactions"
fi

pipenv sync

pipenv run python create.py \
    --channel-seeds-file channel_seed \
    --accounts $AMOUNT_OF_SOURCE_ACCOUNTS \
    --passphrase "$PASS" \
    --horizon $HORIZON \
    --source-account $FUNDER_SEED \
    --json-output True >../../accounts.json

pipenv run python create.py \
    --channel-seeds-file channel_seed \
    --accounts $AMOUNT_OF_DESTINATION_ACCOUNTS \
    --passphrase "$PASS" \
    --horizon $HORIZON \
    --source-account $FUNDER_SEED \
    --json-output True >../../dest.json

# whitelist one address, used to whitelist txs in load test
../../resources/whitelist \
    -horizon $HORIZON \
    -passphrase $PASS \
    -whitelist-seed $WHITELIST_SEED \
    -address-to-whitelist $WHITELISTED_SEED
