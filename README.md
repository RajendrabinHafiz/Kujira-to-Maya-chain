<!-- markdownlint-disable MD041 -->

[![pipeline status](https://gitlab.com/mayachain/mayanode/badges/develop/pipeline.svg)](https://gitlab.com/mayachain/mayanode/commits/develop)
[![coverage report](https://gitlab.com/mayachain/mayanode/badges/develop/coverage.svg)](https://gitlab.com/mayachain/mayanode/-/commits/develop)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/mayachain/mayanode)](https://goreportcard.com/report/gitlab.com/mayachain/mayanode)

# MAYAChain

---

> **Mirror**
>
> This repo mirrors from MAYAChain Gitlab to Github.
> To contribute, please contact the team and commit to the Gitlab repo:
>
> https://gitlab.com/mayachain/mayanode

---

======================================

MAYAChain is a decentralised liquidity network built with [CosmosSDK](cosmos.network).

## MAYANodes

The MAYANode software allows a node to join and service the network, which will run with a minimum of four nodes. The only limitation to the number of nodes that can participate is set by the `minimumBondAmount`, which is the minimum amount of capital required to join. Nodes are not permissioned; any node that can bond the required amount of capital can be scheduled to churn in.

MAYAChain comes to consensus about events observed on external networks via witness transactions from nodes. Swap and liquidity provision logic is then applied to these finalised events. Each event causes a state change in MAYAChain, and some events generate an output transaction which require assets to be moved (outgoing swaps or bond/liquidity withdrawals). These output transactions are then batched, signed by a threshold signature scheme protocol and broadcast back to the respective external network. The final gas fee on the network is then accounted for and the transaction complete.

This is described as a "1-way state peg", where only state enters the system, derived from external networks. There are no pegged tokens or 2-way pegs, because they are not necessary. On-chain Bitcoin can be swapped with on-chain Ethereum in the time it takes to finalise the confirmed event.

All funds in the system are fully accounted for and can be audited. All logic is fully transparent.

## Churn

MAYAChain actively churns its validator set to prevent stagnation and capture, and ensure liveness in signing committees. Churning is also the mechanism by which the MAYANode software can safely facilitate non-contentious upgrades.

Every 50000 blocks (3 days) MAYAChain will schedule the oldest and the most unreliable node to leave, and rotate in two new nodes. The next two nodes chosen are simply the nodes with the highest bond.

During a churn event the following happens:

- The incoming nodes participate in a TSS key-generation event to create new Asgard vault addresses
- When successful, the new vault is tested with a on-chain challenge-response.
- If successful, the vaults are rolled forward, moving all assets from the old vault to the new vault.
- The outgoing nodes are refunded their bond and removed from the system.

## Bifröst

The Bifröst faciliates connections with external networks, such as Binance Chain, Ethereum and Bitcoin. The Bifröst is generally well-abstracted, needing only minor changes between different chains. The Bifröst handles observations of incoming transactions, which are passed into MAYAChain via special witness transactions. The Bifröst also handles multi-party computation to sign outgoing transactions via a Genarro-Goldfeder TSS scheme. Only 2/3rds of nodes are required to be in each signing ceremony on a first-come-first-serve basis, and there is no log of who is present. In this way, each node maintains plausible deniabilty around involvement with every transaction.

### Adding a New Chain

To add a new chain, the process and guidelines defined in [this doc](docs/chains/README.md) must be followed.

### Removing a Chain

To remove a chain, nodes can stop witnessing it. If a super-majority of nodes do not promptly follow suit, the non-witnessing nodes will attract penalties during the time they do not witness it. If a super-majority of nodes stop witnessing a chain it will invoke a chain-specific Ragnörok, where all funds attributed to that chain will be returned and the chain delisted.

## Transactions

The MAYAChain facilitates the following transactions, which are made on external networks and replayed into the MAYAChain via witness transactions:

- **ADD LIQUIDITY**: Anyone can provide assets in pools. If the asset hasn't been seen before, a new pool is created.
- **WITHDRAW LIQUIDITY**: Anyone who is providing liquidity can withdraw their claim on the pool.
- **SWAP**: Anyone can send in assets and swap to another, including sending to a destination address, and including optional price protection.
- **BOND**: Anyone can bond assets and attempt to become a Node. Bonds must be greater than the `minimumBondAmount`, else they will be refunded.
- **LEAVE**: Nodes can voluntarily leave the system and their bond and rewards will be paid out. Leaving takes 6 hours.
- **RESERVE**: Anyone can add assets to the Protocol Reserve, which pays out to Nodes and Liquidity Providers. 220,447,472 Rune will be funded in this way.

## Continuous Liquidity Pools

The Provision of liquidity logic is based on the `CLP` Continuous Liquidity Pool algorithm.

**Swaps**
The algorithm for processing assets swaps is given by:
`y = (x * Y * X) / (x + X)^2`, where `x = input, X = Input Asset, Y = Output Asset, y = output`

The fee paid by the trader is given by:
`fee = ( x^2 * Y ) / ( x + X )^2 `

The slip-based fee model has the following benefits:

- Resistant to manipulation
- A proxy for demand of liquidity
- Asymptotes to zero over time, ensuring pool prices match reference prices
- Prevents Impermanent Loss to liquidity providers

**Provide Liquidity**
The provider units awarded to a liquidity provider is given by:
`liquidityUnits = ((R + T) * (r * T + R * t))/(4 * R * T)`, where `r = Rune Provided, R = Rune Balance, T = Token Balance, t = Token Provided`

This allows them to provide liquidity asymmetrically since it has no opinion on price.

## Incentives

The system is safest and most capital-efficient when 67% of Rune is bonded and 33% is provided liquidity in pools. At this point, nodes will be paid 67% of the System Income, and liquidity providers will be paid 33% of the income. The Sytem Income is the block rewards (`blockReward = totalReserve / 6 / 6311390`) plus the liquidity fees collected in that block.

An Incentive Pendulum ensures that liquidity providers receive 100% of the income when 0% is provided liquidity (inefficent), and 0% of the income when `totalLiquidity >= totalBonded` (unsafe).
The Total Reserve accumulates the `transactionFee`, which pays for outgoing gas fees and stabilises long-term value accrual.

## Governance

There is strictly minimal goverance possible through MAYANode software. Each MAYANode can only generate valid blocks that is fully-compliant with the binary run by the super-majority.

The best way to apply changes to the system is to submit a MAYAChain Improvement Proposal (TIP) for testing, validation and discussion among the MAYAChain developer community. If the change is beneficial to the network, it can be merged into the binary. New nodes may opt to run this updated binary, signalling via a `semver` versioning scheme. Once the super-majority are on the same binary, the system will update automatically. Schema and logic changes can be applied via this approach.

Changes to the Bifröst may not need coordination, as long as the changes don't impact MAYAChain schema or logic, such as adding new chains.

Emergency changes to the protocol may be difficult to coordinate, since there is no ability to communicate with any of the nodes. The best way to handle an emergency is to invoke Ragnarök, simply by leaving the system. When the system falls below 4 nodes all funds are paid out and the system can be shut-down.

======================================

## Setup

Install dependencies, you may skip packages you already have.

```bash
apt-get update
apt-get install -y git make golang-go protobuf-compiler
```

Install [Docker and Docker Compose V2](https://docs.docker.com/engine/install/).

Ensure you have a recent version of go ([scripts/check-build-env.sh](https://gitlab.com/mayachain/mayanode/-/blob/develop/scripts/check-build-env.sh#L7-9)) and enabled go modules.<br/>
Add `GOBIN` to your `PATH`.

```bash
export GOBIN=$GOPATH/bin
```

### Automated Install Locally

Clone repo

```bash
git clone https://gitlab.com/mayachain/mayanode.git
cd thornode
```

Install via this `make` command.

```bash
make openapi
make install
```

Once you've installed `thornode`, check that they are there.

```bash
thornode help
```

### Start Standalone Full Stack

For development and running a full chain locally (your own separate network), use the following command on the project root folder:

```bash
make run-mocknet
```

See [build/docker/README.md](./build/docker/README.md) for more detailed documentation on the MAYANode images and local mocknet environment.

### Simulate Local Churn

```bash
# reset mocknet cluster
make reset-mocknet-cluster

# increase churn interval as desired from the default 60 blocks
make cli-mocknet
> thornode tx mayachain mimir CHURNINTERVAL 1000 --from dog $TX_FLAGS

# bootstrap vaults from smoke test add liquidity transactions
make mocknet-bootstrap

# verify vault balances
curl -s localhost:1317/thorchain/vaults/asgard | jq '.[0].coins'

# watch logs for churn
make logs-mocknet

# verify active nodes
curl -s localhost:1317/thorchain/nodes | jq '[.[]|select(.status=="Active")]|length'

# disable future churns if desired
make cli-mocknet
> thornode tx mayachain mimir CHURNINTERVAL 1000000 --from dog $TX_FLAGS
```

See [build/docker/README.md](./build/docker/README.md) for more detailed documentation on the MAYANode images and local mocknet environment.

### Smoke Tests

The smoke tests compare a mocknet against a simulator implemented in python. Changes to thornode, particularly to the calculations, will require also updating the python simulator, and subsequently the unit-tests for the simulator.

The smoke-test currently requires that all synth balances be cleared be liquidity is withdrawn at the end of the smoke-test, so it is possible the transactions in `test/smoke/data/smoke_test_transactions.json` may need to be changed.

#### Run Smoke Tests

```bash
make smoke-protob-docker
make smoke
```

#### Update Balances and Events

```bash
EXPORT=data/smoke_test_balances.json EXPORT_EVENTS=data/smoke_test_events.json make smoke-unit-test
```

### Format code

```bash
make format
```

### Build all

```bash
make all
```

### Test

Run tests

```bash
make test
```

To run test live when you change a file, use...

```bash
go get -u github.com/mitranim/gow
make test-watch
```

### How to contribute

Check [contributing](./CONTRIBUTING.md)
