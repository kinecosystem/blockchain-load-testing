#!/usr/bin/env python
"""Process load testing log from stdin and print transaction response time in CSV format.

This script is obsolete and is superseded by tx_ledger_time_diff.py
"""
import csv
import fileinput
import json
import logging
import sys


logging.basicConfig(level=logging.DEBUG)


def response_time_csv():
    """Read load testing log from stdin and print transaction response time in CSV format."""
    # convert nanosecond to millisecond with correct rounding
    # e.g. (1.1 --> 1) and (1.9 --> 2)
    #
    # https://stackoverflow.com/questions/3950372/round-with-integer-division
    ns_to_ms = lambda x: (x + 1000000 // 2) // 1000000

    # generate response time bucket dictionary
    # for sorting transactions by buckets of (0ms, 100ms, 200ms, ... , 5 min)
    buckets = {x*100: 0 for x in range(3000)}

    # count response time for successful transactions
    # and assign to response time buckets
    for i, raw_line in enumerate(fileinput.input()):
        if i % 100 == 0:
            logging.debug('processing line %d', i)

        json_line = json.loads(raw_line)

        if json_line.get('transaction_status') == 'success':
            response_time = ns_to_ms(json_line['response_time_nanoseconds'])
            # zero out two last digits e.g. 125 --> 100
            buckets[response_time // 100 * 100] += 1

    w = csv.writer(sys.stdout)
    w.writerow(['response time (ms)', 'count'])

    for b in buckets.items():
        w.writerow(b)


if __name__ == '__main__':
    response_time_csv()
