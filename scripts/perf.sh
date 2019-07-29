#!/bin/bash

. vars.sh

if [ -e loadtest.log ]; then rm loadtest.log; fi

#rate
if [ -z "$1" ]
then
      echo "\$1 is empty, using value from vars.sh"
else
      export RATE=$1
fi

#Repetitions
if [ -z "$2" ]
then
      echo "\$2 is empty, using value from vars.sh"
else
      export REPETITIONS=$2
fi



for (( i=1; i<=$REPETITIONS; i++ )); do

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

echo "Gathering logs."

./scripts/perf-logs.sh $L1 $L2

mv loadtest.log /tmp/loadtest-$L1-$L2.log

SCRIPT_DIR=`pwd`/scripts

pushd /tmp

a=$(wc -l perf-tx-ledgers.txt | cut -d" " -f 1)
cat perf-tx-ledgers.txt | tail -n $(($a - 2)) | head -n $(($a - 2 - 2)) > perf-tx-ledgers.txt2
mv perf-tx-ledgers.txt2 perf-tx-ledgers.txt

awk -F'|' -f $SCRIPT_DIR/tx_ledger.awk perf-tx-ledgers.txt > tx-ledgers.sql
cat tx-ledgers.sql | ssh perf-horizon 'psql -h perf-horizon-1.c3ofnhfeeiwz.us-west-1.rds.amazonaws.com -U stellar analytics' > /dev/null 2>&1

cat perf-horizon-ingest.log | . $SCRIPT_DIR/ingestion.sh > ingestion.sql
cat ingestion.sql | ssh perf-horizon 'psql -h perf-horizon-1.c3ofnhfeeiwz.us-west-1.rds.amazonaws.com -U stellar analytics' > /dev/null 2>&1

grep "submitting" loadtest-$L1-$L2.log | jq -rc ".tx_hash, .timestamp" | awk -f $SCRIPT_DIR/submission.awk | paste -sd ",\n" | while read -r line; do echo "insert into submission values($line);"; done > submission.sql
cat submission.sql | ssh perf-horizon 'psql -h perf-horizon-1.c3ofnhfeeiwz.us-west-1.rds.amazonaws.com -U stellar analytics' > /dev/null 2>&1

if [ -e core.sql ]; then rm core.sql; fi
for file in `ls perf*$L1-$L2-log.json.gz`; do
	$SCRIPT_DIR/core_stats.sh $file >> core.sql
done
cat core.sql | ssh perf-horizon 'psql -h perf-horizon-1.c3ofnhfeeiwz.us-west-1.rds.amazonaws.com -U stellar analytics' > /dev/null 2>&1

gzip loadtest-$L1-$L2.log
gzip -f perf-tx-ledgers.txt
gzip -f perf-horizon-ingest.log
gzip -f tx-ledgers.sql
gzip -f ingestion.sql
gzip -f submission.sql
gzip -f core.sql

TAR=logs-$L1-$L2.$TIME_LENGTH.$RATE.horizon.tar
tar cvf $TAR perf*$L1-$L2* loadtest-$L1-$L2.log.gz test-params perf-tx-ledgers.txt.gz perf-horizon-ingest.log.gz *.sql.gz 
aws s3api --no-sign-request put-object --bucket perf-test-s3-logs --key $TAR --body $TAR
popd

done
