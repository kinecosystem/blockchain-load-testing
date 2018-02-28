#!/usr/bin/env python
"""Process load testing log from stdin and print transaction success rate in CSV format."""
import csv
import fileinput
import json
import logging
import sys


logging.basicConfig(level=logging.DEBUG)


def success_rate():
    """Process load testing log from stdin and print transaction success rate in CSV format."""
    success, failure = 0, 0

    for i, raw_line in enumerate(fileinput.input()):
        if i % 100 == 0:
            logging.debug('processing line %d', i)

        json_line = json.loads(raw_line)

        # get transaction status from log event
        # NOTE not all lines are even transaction status so this field might be
        # missing, in which case we skip that line
        try:
            status = json_line['transaction_status']
        except KeyError:
            continue

        if status == 'success':
            success += 1
        elif status == 'failure':
            failure += 1
        else:
            raise RuntimeError('unknown transaction status')

    w = csv.writer(sys.stdout)
    w.writerow(['type', 'count'])

    w.writerow(['success', success])
    w.writerow(['failure', failure])


if __name__ == '__main__':
    success_rate()
