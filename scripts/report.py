#!/usr/bin/env python

import collections
import csv
import datetime
import json
import sys

RFC3339_NO_TIMEZONE = '%Y-%m-%dT%H:%M:%S'
RFC3339 = '%Y-%m-%dT%H:%M:%S.%f'  # TODO output burst limit


def response_time_csv(path, log):
    # convert nanosecond to millisecond with correct rounding
    # e.g. (1.1 --> 1) and (1.9 --> 2)
    #
    # https://stackoverflow.com/questions/3950372/round-with-integer-division
    ns_to_ms = lambda x: (x + 1000000 // 2) // 1000000

    # generate response time bucket dictionary
    # for sorting transactions by buckets of (0ms, 100ms, 200ms, ... , 1 min)
    buckets = {x*100: 0 for x in range(600)}

    for line in log:
        if line.get('status') == 'success':
            response_time = ns_to_ms(line['response_time_nanoseconds'])
            # zero two last digits e.g. 125 --> 100
            buckets[response_time // 100 * 100] += 1


    with open(path, 'w') as f:
        w = csv.writer(f)
        w.writerow(['buckets', 'response time (ms)'])

        for b in buckets.items():
            w.writerow(b)


def success_rate(path, log):
    success, failure = 0, 0

    for line in log:
        status = line.get('status')
        if status == 'success':
            success += 1
        elif  status == 'failure':
            failure += 1

    with open(path, 'w') as f:
        w = csv.writer(f)
        w.writerow(['type', 'count'])

        w.writerow(['success', success])
        w.writerow(['failure', failure])

def tx_rate(path, log):
    buckets = collections.OrderedDict()
    for line in log:
        if line.get('msg') == 'submitting transaction':
            ts_str = line['timestamp']
            i = ts_str.rfind('.')

            ts = datetime.datetime.strptime(ts_str[:i], RFC3339_NO_TIMEZONE)

            unix_ts = str(int(ts.timestamp()))
            if unix_ts not in buckets:
                buckets[unix_ts] = 1
            else:
                buckets[unix_ts] += 1

    with open(path, 'w') as f:
        w = csv.writer(f)
        w.writerow(['timestamp', 'transactions'])

        for b in buckets.items():
            w.writerow(b)


# TODO
# def tx_burst_rate(path, log):
#     buckets = collections.OrderedDict()
#     for line in log:
#         if line.get('msg') == 'submitting transaction':
#             ts_str = line['timestamp']
#             i = ts_str.rfind('.')

#             ts = datetime.datetime.strptime(ts_str[:i], RFC3339)

#             unix_ts = str(int(ts.timestamp()))
#             if unix_ts not in buckets:
#                 buckets[unix_ts] = 1
#             else:
#                 buckets[unix_ts] += 1

#     with open(path, 'w') as f:
#         w = csv.writer(f)
#         w.writerow(['timestamp', 'transactions'])

#         for b in buckets.items():
#             w.writerow(b)


def read_log(path):
    log = []
    with open(path, 'r') as f:
        for l in f:
            log.append(json.loads(l))

    return log


if __name__ == '__main__':
    log = read_log(sys.argv[1])
    # response_time_csv('response_times.csv', log)
    # success_rate('success_rate.csv', log)
    tx_rate('tx_rate.csv', log)
