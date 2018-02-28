#!/usr/bin/env python
"""Load all transactions information for given account and print it."""
import json
import logging
import sys

import requests


HORIZON = 'https://horizon.stellar.org'
ORDER = 'asc'
LIMIT = 200  # maximum amount of reuslts per request

RECORD_FIELDS = ("id", "hash", "ledger", "source_account", "created_at")


logging.basicConfig(level=logging.DEBUG)


def print_txs(horizon, address):
    """Load all transaction information for given address from Horizon and print."""
    params = {'order': ORDER, 'limit': LIMIT}
    c = 0
    while True:
        # set paging token
        try:
            params['cursor'] = records[-1]['paging_token']
        except NameError:
            # exception should be raised (an ignored) on first iteration only
            # when the paging token shouldn't be used
            pass

        # request transaction information
        logging.debug('requesting transactions, iteration %d', c)
        r = requests.get('{horizon}/accounts/{address}/transactions'.format(horizon=horizon, address=address),
                         params=params)
        r.raise_for_status()
        records = r.json()['_embedded']['records']

        # stop paging when no more records are returned
        if not records:
            break

        for record in records:
            # print subset of fields, since the object is very large
            print(json.dumps({k: record[k] for k in RECORD_FIELDS}))

        c += 1

if __name__ == '__main__':
    address = sys.argv[1]
    print_txs(HORIZON, address)
