#!/bin/sh
#IMPORTANNT: FUNDER_SEED and CHANNEL_SEED will be overridden if activated from Jenkins

#needs to be lower than AMOUNT_OF_SOURCE_ACCOUNTS
export RATE=100
export REPETITIONS=1
export TIME_LENGTH=300

#max amount of simultaneous tx
export AMOUNT_OF_SOURCE_ACCOUNTS=1000
export ACCOUNTS_FILE=accounts.json
export GOPATH=~/go
export AMOUNT_OF_DESTINATION_ACCOUNTS=20
export DEST_ACCOUNT=dest.json

export OPS_PER_TX=10
#less than or equal to the AMOUNT_OF_SOURCE_ACCOUNTS
export SUBMITTERS=500

export SCRIPT_DIR=~/stellar-load-testing/scripts


export CORE_SERVERS="ip-core-test-1.test.kin ip-core-test-2.test.kin  ip-core-test-3.test.kin ip-core-test-4.test.kin ip-core-test-5.test.kin"
export HORIZON=http://ip-horizon-test-1.test.kin
export PROMETHEUS=ip-prometheus_server-test-1.test.kin

export S3_BUCKET="perf-test-s3-logs"

export PASSPHRASE="Kin test ; Jul 2019"

#test seeds
export FUNDER_SEED=SCEIO4XV4UFZEVFAWOWCSXUA32LKLKG5MLJ5R7FVCGVSLFTIHP3A35E7
export CHANNEL_SEED=SCEIO4XV4UFZEVFAWOWCSXUA32LKLKG5MLJ5R7FVCGVSLFTIHP3A35E7



