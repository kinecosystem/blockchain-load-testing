"""Create accounts using given channel seeds (as sequence number consumers) and root account seed (as source account)."""
import asyncio
import argparse
import json
import logging

from kin import KinClient, Environment, Keypair
from kin.blockchain.builder import Builder

from helpers import (NETWORK_NAME, MIN_FEE,
                     load_accounts, generate_keypairs, get_sequences_multiple_endpoints,
                     create_accounts)


STARTING_BALANCE = 1e5


logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s %(message)s')


def parse_args():
    """Generate and parse CLI arguments."""
    parser = argparse.ArgumentParser()

    parser.add_argument('--source-account', required=True, type=str, help='Source account to pay transaction fees')
    parser.add_argument('--channel-seeds-file', required=True, type=str, help='File path to channel seeds file')
    parser.add_argument('--accounts', required=True, type=int, help='Amount of accounts to create')
    parser.add_argument('--passphrase', required=True, type=str, help='Network passphrase')
    parser.add_argument('--horizon', action='append', help='Horizon endpoint URL (use multiple --horizon flags for multiple addresses)')
    parser.add_argument('--json-output', required=False, type=bool, help='Export output to json format')
    return parser.parse_args()


async def init_channel_builders(channel_seeds_file, passphrase, horizon):
    env = Environment(NETWORK_NAME, (horizon)[0], passphrase)

    channel_kps = load_accounts(channel_seeds_file)
    sequences = await get_sequences_multiple_endpoints(horizon, [kp.public_address for kp in channel_kps])

    channel_builders = []
    for i, kp in enumerate(channel_kps):
        b = Builder(NETWORK_NAME, KinClient(env).horizon, MIN_FEE, kp.secret_seed)
        b.sequence = str(sequences[i])
        channel_builders.append(b)

    return channel_builders


async def main():
    """Create accounts and print their seeds to stdout."""
    args = parse_args()

    # initialize channels
    channel_builders = await init_channel_builders(args.channel_seeds_file, args.passphrase, args.horizon)
    kps = generate_keypairs(args.accounts)
    source_kp = Keypair(args.source_account)

    await create_accounts(source_kp, kps, channel_builders, args.horizon, STARTING_BALANCE)

    if args.json_output:
        keypairs = []
        for kp in kps:
            keypairs.append({"address": kp.public_address, "seed": kp.secret_seed})

        out = {"keypairs": keypairs}
        print(json.dumps(out, indent=True))
    else:
        out = (kp.secret_seed for kp in kps)
        print('\n'.join(list(out)))


if __name__ == '__main__':
    asyncio.run(main())
