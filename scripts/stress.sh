#!/bin/bash

rm loadtest.log;

for (( i=1; i<=$2; i++)); do
	./scripts/loadtest.sh $1 true > /tmp/a && \
		sleep 15 && \
		curl -sS "$HORIZON/ledgers?limit=200&order=desc" > /tmp/$1 && \
		PT=`grep paging_token /tmp/$1 | tail -1 | egrep -o '([0-9]*)'` && \
		curl -sS "$HORIZON/ledgers?limit=200&order=desc&cursor=$PT" >> /tmp/$1
done
