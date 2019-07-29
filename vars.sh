#!/bin/sh
#IMPORTANNT: FUNDER_SEED and CHANNEL_SEED will be overridden if activated from Jenkins

export RATE=80
export REPETITIONS=1

export AMOUNT_OF_SOURCE_ACCOUNTS=10
export ACCOUNTS_FILE=accounts.json
export GOPATH=~/go
export AMOUNT_OF_DESTINATION_ACCOUNTS=2
export DEST_ACCOUNT=dest.json

export OPS_PER_TX=1
export SUBMITTERS=2000

CORE_SERVERS="core-test-1.test.kin core-test-2.test.kin  core-test-3.test.kin core-test-4.test.kin core-test-5.test.kin"
S3_BUCKET_NAME="perf-test-s3-logs"

#CORE_SERVERS="perf-1 perf-2 perf-3 perf-4 perf-5 perf-6 perf-7 perf-8 perf-9 perf-10 perf-11"
export PASSPHRASE="Kin test ; Jul 2019"
export HORIZON=http://ip-horizon-test-1.test.kin

export FUNDER_SEED=xxx
export CHANNEL_SEED=yyy

export FUNDER_SEED=SCEIO4XV4UFZEVFAWOWCSXUA32LKLKG5MLJ5R7FVCGVSLFTIHP3A35E7
export CHANNEL_SEED=SCEIO4XV4UFZEVFAWOWCSXUA32LKLKG5MLJ5R7FVCGVSLFTIHP3A35E7
