#!/bin/sh
#IMPORTANNT: FUNDER_SEED and CHANNEL_SEED will be overridden if activated from Jenkins

export RATE=80
export REPETITIONS=1
export TIME_LENGTH=600

export AMOUNT_OF_SOURCE_ACCOUNTS=10
export ACCOUNTS_FILE=accounts.json
export GOPATH=~/go
export AMOUNT_OF_DESTINATION_ACCOUNTS=2
export DEST_ACCOUNT=dest.json

export OPS_PER_TX=1
export SUBMITTERS=2000

export CORE_SERVERS="ip-core-test-1.test.kin ip-core-test-2.test.kin  ip-core-test-3.test.kin ip-core-test-4.test.kin ip-core-test-5.test.kin"
export HORIZON=http://ip-horizon-test-1.test.kin


export S3_BUCKET_NAME="perf-test-s3-logs"

#CORE_SERVERS="p    erf-1 perf-2 perf-3 perf-4 perf-5 perf-6 perf-7 perf-8 perf-9 perf-10 perf-11"
export PASSPHRASE="Kin test ; Jul 2019"

#test seeds
export FUNDER_SEED=SCEIO4XV4UFZEVFAWOWCSXUA32LKLKG5MLJ5R7FVCGVSLFTIHP3A35E7
export CHANNEL_SEED=SCEIO4XV4UFZEVFAWOWCSXUA32LKLKG5MLJ5R7FVCGVSLFTIHP3A35E7

export SCRIPT_DIR=~/stellar-load-testing/scripts
