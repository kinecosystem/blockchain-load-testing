# TODO split these into module files with explicit names other than generic "helpers"
"""Common helper functions for tests."""
import asyncio
import concurrent.futures
import logging
import multiprocessing
import math
from hashlib import sha256
from typing import List

import aiohttp

from kin import Keypair
from kin.blockchain.builder import Builder
from kin_base import Keypair as BaseKeypair

NETWORK_NAME = 'LOCAL'
MIN_FEE = 100
MAX_OPS = 100
TX_SET_SIZE = 500


def root_account_seed(passphrase: str) -> str:
    """Return the root account seed based on the given network passphrase."""
    network_hash = sha256(passphrase.encode()).digest()
    return BaseKeypair.from_raw_seed(network_hash).seed().decode()


def derive_root_account(passphrase):
    """Return the keypair of the root account, based on the network passphrase."""
    network_hash = sha256(passphrase.encode()).digest()
    seed = BaseKeypair.from_raw_seed(network_hash).seed().decode()
    return Keypair(seed)


# XXX concurrent.futures behavior: why can't this function be a local function in generate_keypairs()?
def keypair_list(n) -> List[Keypair]:
    """Return Keypair list according to given amount."""
    return [Keypair() for _ in range(n)]


def load_accounts(path) -> List[Keypair]:
    """Load seeds from file path and return Keypair list.

    Expected file format is a newline-delimited seed list.
    """
    kps = []
    with open(path) as f:
        for seed in f:
            kps.append(Keypair(seed.strip()))
    return kps


def generate_keypairs(n) -> List[Keypair]:
    """Generate Keypairs efficiently using all available CPUs."""
    logging.info('generating %d keypairs', n)

    # split amounts of keypairs to create to multiple inputs,
    # one for each cpu
    cpus = multiprocessing.cpu_count()
    d, m = n // cpus, n % cpus  # d[iv], m[od]
    keypair_amounts = [d]*cpus + [m]*(1 if m else 0)

    # generate keypairs across multiple cpus
    with concurrent.futures.ProcessPoolExecutor() as process_pool:
        futurs = []
        for amount in keypair_amounts:
            f = process_pool.submit(keypair_list, amount)
            futurs.append(f)

    kp_batches, _ = concurrent.futures.wait(futurs)
    kps = []
    for batch in kp_batches:
        kps.extend(batch.result())

    logging.info('%d keypairs generated', n)
    return kps


def _build_and_sign(source_kp, channel, kps, starting_balance):
    for kp in kps:
        channel.append_create_account_op(source=source_kp.public_address,
                                         destination=kp.public_address,
                                         starting_balance=str(starting_balance))

    channel_seed = channel.keypair.seed().decode()
    channel.sign(secret=channel_seed)
    if channel_seed != source_kp.secret_seed:
        channel.sign(secret=source_kp.secret_seed)

    xdr = channel.gen_xdr().decode()
    channel.clear()

    return xdr


async def channel_create_accounts(pool, session, queue, source_kp: Keypair, kps: List[Keypair], starting_balance, horizon_endpoint):
    """Create MAX_OPS accounts in a single transaction using given channel and horizon endpoint."""
    # get avaiable channel i.e. one which isn't currently in the process of submitting a tx
    # and use it to async submit a single create account tx.
    channel = await queue.get()
    logging.debug('using channel %s to create accounts', channel)

    # sign transaction with channel and source accounts,
    # utilizing all availables cpus (tx signing is computationally intensive)
    loop = asyncio.get_running_loop()
    xdr = await loop.run_in_executor(pool, _build_and_sign, source_kp, channel, kps, starting_balance)

    # submit tx
    await post(session, '{}/transactions'.format(horizon_endpoint), {'tx': xdr}, [200])

    # update sequence number
    res = await get(
        session,
        '{}/accounts/{}'.format(horizon_endpoint, channel.keypair.address().decode()),
        None,
        [200])

    channel.sequence = int(res['sequence'])

    # make channel available for another transaction
    await queue.put(channel)
    logging.debug('channel %s finished submitting current create account transaction', channel)


async def create_accounts(source_kp: Keypair, account_kps: List[Keypair], channel_builders: List[Builder], horizon_endpoints: List[str], starting_balance: int):
    """Asynchronously create accounts and return a Keypair instance for each created account.

    Accounts are created using given channel builders (as sequence number consumers)
    and source keypair as funding source account.
    """
    logging.info('creating %d accounts', len(account_kps))

    # generate txs, squeezing as much "create account" ops as possible to each one.
    # when each tx is full with as much ops as it can include, sign and generate
    # that tx's XDR.
    # then, continue creating accounts using a new tx, and so on.
    # we stop when we create all ops required according to given accounts_num.
    def batch(iterable, n=1):
        l = len(iterable)
        for ndx in range(0, l, n):
            yield iterable[ndx:min(ndx + n, l)]

    # put channels in queue.
    # each channel will be used to submit create account txs asynchronously,
    # and will continue submitting more txs after he's done with the previous
    # txs
    channel_queue = asyncio.Queue()
    for c in channel_builders:
        await channel_queue.put(c)

    with concurrent.futures.ProcessPoolExecutor() as process_pool:
        async with LoggingClientSession(connector=aiohttp.TCPConnector(limit=len(channel_builders))) as session:
            futurs = []
            for i, kp_batch in enumerate(batch(account_kps, MAX_OPS)):
                coro = channel_create_accounts(
                    process_pool, session, channel_queue, source_kp, kp_batch, starting_balance, horizon_endpoints[i % len(horizon_endpoints)])

                futurs.append(asyncio.create_task(coro))

            # wait for all remaining transactions to finish
            await asyncio.gather(*futurs)

    # sanity check: since all txs have already finished, all channels should be
    # back in the queue
    assert channel_queue.qsize() == len(channel_builders)

    logging.info('created %d accounts', len(account_kps))


async def get(session: aiohttp.ClientSession, url, req_params, expected_statuses: List[int]):
    """Send an HTTP GET request and return response JSON data.

    Fail if response isn't expected status code or format other than JSON.
    """
    async with session.get(url, params=req_params) as res:
        try:
            res_data = await res.json()
        except aiohttp.client_exceptions.ContentTypeError as e:
            logging.error(e)
            logging.error(await res.text())

        if res.status not in expected_statuses:
            logging.error('Error in HTTP GET request to %s: %s', url, res_data)
            raise RuntimeError('Error in HTTP GET request to {}'.format(url))

    return res_data


async def post(session: aiohttp.ClientSession, url, req_data, expected_statuses: List[int]):
    """Send an HTTP POST request with given data and return response JSON data.

    Fail if response isn't expected status code or format other than JSON.

    NOTE expected status is either OK 200 or Server Timeout 504 if the transaction
    wasn't added to the next three ledgers.
    """
    async with session.post(url, data=req_data) as res:
        try:
            res_data = await res.json()
        except aiohttp.client_exceptions.ContentTypeError as e:
            logging.error(e)
            logging.error(await res.text())

        if res.status not in expected_statuses:
            logging.error('Error in HTTP POST request to %s with data %s: %s', url, req_data, res_data)
            raise RuntimeError('Error in HTTP POST request to {}'.format(url))

    return res_data


class LoggingClientSession(aiohttp.ClientSession):
    """aiohttp client session that logs requests."""

    async def _request(self, method, url, **kwargs):
        logging.debug('Starting request <%s %r>', method, url)
        return await super()._request(method, url, **kwargs)


async def send_txs_multiple_endpoints(endpoints, xdrs, submit_to_horizon=True, expected_statuses=(200, 504)):
    """Send multiple async transaction XDRs submitting to one of given endpoints.

    endpoints are iterated one after the other in a round robin manner.

    Can submit both to Horizon (HTTP POST) or Core (HTTP GET).
    """
    logging.info('sending %d transactions', len(xdrs))

    async with LoggingClientSession(connector=aiohttp.TCPConnector(limit=5000)) as session:
        # generate urls
        if submit_to_horizon:
            urls = ['{}/transactions'.format(e) for e in endpoints]
        else:  # submit to core
            urls = ['{}/tx'.format(e) for e in endpoints]

        # submit to one of the urls in a round robin manner
        results = []
        for i, xdr in enumerate(xdrs):
            url = urls[i % len(urls)]

            if submit_to_horizon:
                coro = post(session, url, {'tx': xdr}, expected_statuses)
            else:  # submit to core
                coro = get(session, url, {'blob': xdr}, expected_statuses)

            results.append(coro)

        results = await asyncio.gather(*results)

    logging.info('%d transactions sent', len(xdrs))
    return results


# TODO remove this function and use send_txs_multiple_endpoints
async def send_txs(txs):
    """Send transactions asynchronously and return responses."""
    with concurrent.futures.ThreadPoolExecutor() as pool:
        loop = asyncio.get_running_loop()
        futurs = [loop.run_in_executor(pool, tx.submit) for tx in txs]
        for _ in await asyncio.gather(*futurs):
            return futurs


async def get_sequences_multiple_endpoints(endpoints, addresses):
    """Get sequence for multiple accounts, using one of given endpoints.

    endpoints are iterated one after the other in a round robin manner.
    """
    logging.info('getting sequence for %d accounts', len(addresses))

    async with LoggingClientSession() as session:
        urls = ['{}/accounts'.format(e) for e in endpoints]
        results = []
        for i, address in enumerate(addresses):
            # send request to one of the urls in a round robin manner
            url = urls[i % len(urls)]
            coro = get(session, '{}/{}'.format(url, address), None, [200])
            results.append(coro)

        results = await asyncio.gather(*results)

    logging.info('finished getting sequence for %d accounts', len(addresses))

    sequences = []
    for r in results:
        try:
            seq = int(r['sequence'])
        except KeyError:
            # can occur if request failed
            seq = 0
        sequences.append(seq)

    return sequences


def get_latest_ledger(client):
    """Return latest ledger dictionary using given KinClient."""
    params = {'order': 'desc', 'limit': 1}
    return client.horizon.ledgers(params=params)['_embedded']['records'][0]


def add_prioritizers(builder: Builder, kps: List[Keypair]):
    """Add given addresses to whitelist account, making them transaction prioritizers."""
    logging.info('adding %d prioritizers', len(kps))

    for batch_index in range(max(1, math.ceil(len(kps)/MAX_OPS))):
        start = batch_index*MAX_OPS
        end = min((batch_index+1)*MAX_OPS,
                  len(kps))

        for i, kp in enumerate(kps[start:end], start=1):
            logging.debug('adding manage data op #%d', i)
            builder.append_manage_data_op(kp.public_address, kp._hint)

        builder.sign()

        logging.debug('submitting transaction with %d manage data ops', end-start)
        builder.submit()
        logging.debug('done')

        builder.clear()

    logging.info('%d prioritizers added', len(kps))
