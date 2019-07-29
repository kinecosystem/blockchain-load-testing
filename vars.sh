#!/bin/sh
#IMPORTANNT: FUNDER_SEED and CHANNEL_SEED will be overridden if activated from Jenkins

export RATE=80
export REPETITION=1

export AMOUNT_OF_SOURCE_ACCOUNTS=10
export ACCOUNTS_FILE=accounts.json
export GOPATH=~/go
export AMOUNT_OF_DESTINATION_ACCOUNTS=2
export DEST_ACCOUNT=dest.json

export OPS_PER_TX=1
export SUBMITTERS=2000


export PASSPHRASE="Kin test ; Jul 2019"
export HORIZON=https://horizon-testnet.kin.org/

#CORE_SERVERS="perf-1 perf-2 perf-3 perf-4 perf-5 perf-6 perf-7 perf-8 perf-9 perf-10 perf-11"
CORE_SERVERS="ip-10-100-112-248.ec2.internal ip-10-100-116-50.ec2.internal ip-10-100-116-96.ec2.internal ip-10-100-116-205.ec2.internal"
S3_BUCKET_NAME="perf-test-s3-logs"

export FUNDER_SEED=xxx
export CHANNEL_SEED=yyy
