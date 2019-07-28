#!/bin/bash

. vars.sh

if [ -e loadtest.log ]; then rm loadtest.log; fi

export RATE=$1

for (( i=1; i<=$2; i++ )); do

cat <<EOF > /tmp/test-params
ops/tx: $OPS_PER_TX
duration: $TIME_LENGTH
rate: $RATE
EOF

L1=`curl -s "$HORIZON/ledgers?limit=1&order=desc" | grep -oP 'sequence": \K\d+'`

echo -n "Running test on: `date`"

./scripts/loadtest.sh $1 true > /tmp/a

L2=`curl -s "$HORIZON/ledgers?limit=1&order=desc" | grep -oP 'sequence": \K\d+'`

echo -n "Finished at: `date`"

./scripts/collect_results.sh

done
