#!/usr/bin/env python
"""Process load testing log from stdin and print transaction submit rate in CSV format.

Submit rate is a table of transaction rate submitted per second
and and a count of one-second-long "time frames" that had this rate. For example:

time 00:00:00 had 8 transactions submitted
time 00:00:01 had 7 transactions submitted
time 00:00:02 had 8 transactions submitted
time 00:00:03 had 8 transactions submitted

The output will be:
7: 1
8: 3

Meaning only a single time window of one second had a rate of 7 transactions,
and 3 time windows of one second had a rate of 8 transactions.
"""
import collections
import csv
import fileinput
import json
import logging
import sys

import strict_rfc3339


logging.basicConfig(level=logging.DEBUG)


def tx_submit_rate():
    """Process load testing log from stdin and print transaction submit rate in CSV format."""
    buckets = collections.OrderedDict()

    for c, raw_line in enumerate(fileinput.input()):
        if c % 100 == 0:
            logging.debug('processing line %d', c)

        json_line = json.loads(raw_line)
        if json_line.get('msg') == 'submitting transaction':
            unix_ts = str(int(strict_rfc3339.rfc3339_to_timestamp(json_line['timestamp'])))
            if unix_ts not in buckets:
                buckets[unix_ts] = 1
            else:
                buckets[unix_ts] += 1

    rate_count = collections.defaultdict(lambda: 0)
    for rate in buckets.values():
        rate_count[rate] += 1


    w = csv.writer(sys.stdout)
    w.writerow(['txs per second (rate, 1s)', 'count'])

    for b in rate_count.items():
        w.writerow(b)


if __name__ == '__main__':
    tx_submit_rate()
