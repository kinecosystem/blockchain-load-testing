#!/bin/bash

. ./vars.sh

for core in $CORE_SERVERS;
do
    echo $core
    ssh -o "StrictHostKeyChecking=no" -i $SSH_KEY ubuntu@$core "cd /data/ ; sudo su -c 'docker-compose logs --no-color  stellar-core | gzip > /tmp/$core-log.json.gz'" && \
    scp -o "StrictHostKeyChecking=no" -i $SSH_KEY ubuntu@$core:/tmp/$core-log.json.gz /tmp/$core-$1-$2-log.json.gz &
done

# give horizon time to commit ???
echo "sleeping - 1 min"
sleep 60
#Extract Horizon domain from vars.sh>Horizon URL
HORIZON_DOMAIN=$(echo "$HORIZON" | awk -F/ '{print $3}')
echo "getting data from $HORIZON_DOMAIN"
ssh  -o "StrictHostKeyChecking=no"  -i $SSH_KEY ubuntu@$HORIZON_DOMAIN  "sudo docker exec data_horizon-db_1 psql -h localhost -U stellar horizon -c 'select transaction_hash, ledger_sequence from history_transactions'" 1> /tmp/perf-tx-ledgers.txt
ssh -i $SSH_KEY ubuntu@$HORIZON_DOMAIN  "sudo su -c 'docker logs horizon   2> /tmp/a'; grep \"Finished ingesting ledgers\" /tmp/a "  | tail -300 > /tmp/perf-horizon-ingest.log



wait
exit 0
