#!/bin/bash


for i in $CORE_SERVERS;
do
    echo $i
    ssh -i $SSH_KEY $i "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/$i-log.json.gz'" && \
    scp -i $SSH_KEY $i:/tmp/$i-log.json.gz /tmp/$i-$1-$2-log.json.gz &
done

#
#ssh perf-1 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf1-log.json.gz'" && \
#scp perf-1:/tmp/perf1-log.json.gz /tmp/perf1-$1-$2-log.json.gz &
#
#ssh perf-2 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf2-log.json.gz'" && \
#scp perf-2:/tmp/perf2-log.json.gz /tmp/perf2-$1-$2-log.json.gz &
#
#ssh perf-3 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf3-log.json.gz'" && \
#scp perf-3:/tmp/perf3-log.json.gz /tmp/perf3-$1-$2-log.json.gz &
#
#ssh perf-4 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf4-log.json.gz'" && \
#scp perf-4:/tmp/perf4-log.json.gz /tmp/perf4-$1-$2-log.json.gz &
#
#ssh perf-5 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf5-log.json.gz'" && \
#scp perf-5:/tmp/perf5-log.json.gz /tmp/perf5-$1-$2-log.json.gz &
#
#ssh perf-6 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf6-log.json.gz'" && \
#scp perf-6:/tmp/perf6-log.json.gz /tmp/perf6-$1-$2-log.json.gz &
#
#ssh perf-7 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf7-log.json.gz'" && \
#scp perf-7:/tmp/perf7-log.json.gz /tmp/perf7-$1-$2-log.json.gz &
#
#ssh perf-8 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf8-log.json.gz'" && \
#scp perf-8:/tmp/perf8-log.json.gz /tmp/perf8-$1-$2-log.json.gz &
#
#ssh perf-9 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf9-log.json.gz'" && \
#scp perf-9:/tmp/perf9-log.json.gz /tmp/perf9-$1-$2-log.json.gz &
#
#ssh perf-10 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf10-log.json.gz'" && \
#scp perf-10:/tmp/perf10-log.json.gz /tmp/perf10-$1-$2-log.json.gz &
#
#ssh perf-11 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf11-log.json.gz'" && \
#scp perf-11:/tmp/perf11-log.json.gz /tmp/perf11-$1-$2-log.json.gz &

#ssh perf-12 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf12-log.json.gz'" && \
#scp perf-12:/tmp/perf12-log.json.gz /tmp/perf12-$1-$2-log.json.gz &

#ssh perf-13 "sudo su -c 'docker-compose logs --no-color --tail=300000 | gzip > /tmp/perf13-log.json.gz'" && \
#scp perf-13:/tmp/perf13-log.json.gz /tmp/perf13-$1-$2-log.json.gz &

# give horizon time to commit ???
#sleep 60
#ssh perf-horizon "psql -h perf-horizon-1.c3ofnhfeeiwz.us-west-1.rds.amazonaws.com -U stellar horizon -c 'select transaction_hash, ledger_sequence from history_transactions where ledger_sequence >= $1 and ledger_sequence <= $2;'" 1> /tmp/perf-tx-ledgers.txt
#ssh perf-horizon ". log.sh" | tail -300 > /tmp/perf-horizon-ingest.log

wait
exit 0
