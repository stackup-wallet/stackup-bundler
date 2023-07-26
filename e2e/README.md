# Bundler E2E tests

A repeatable set of E2E tests to automate QA checks for the bundler. This should be used in addition to the [bundler test suite](https://github.com/eth-infinitism/bundler-spec-tests).

# Usage

Below are instructions on how to run a series of E2E tests to check that everything is working as expected. The tests will execute a collection of known transactions that cover a wide range of edge cases.

## Prerequisites

The steps in the following section assumes that all these tools have been installed and ready to go.

- Node.JS >= 18
- [Geth](https://geth.ethereum.org/docs/getting-started/installing-geth)

## Setting the environment

To reduce the impact of external factors, we'll run the E2E test using an isolated local instance of both geth and the bundler.

First, we'll need to run a local instance of geth with the following command:

```bash
geth \
  --http.vhosts '*,localhost,host.docker.internal' \
  --http \
  --http.api eth,net,web3,debug \
  --http.corsdomain '*' \
  --http.addr "0.0.0.0" \
  --nodiscover --maxpeers 0 --mine \
  --networkid 1337 \
  --dev \
  --allow-insecure-unlock \
  --rpc.allow-unprotected-txs \
  --miner.gaslimit 12000000
```

In a separate process, navigate to the [eth-infinitism/account-abstraction](https://github.com/eth-infinitism/account-abstraction/) directory and run the following command to deploy the required contracts:

```bash
yarn deploy --network localhost
```

Next, navigate to the [stackup-wallet/contracts](https://github.com/stackup-wallet/contracts) directory and run the following command to deploy the supporting test contracts:

```bash
yarn deploy:AllTest --network localhost
```

Lastly, run the bundler with the following config:

```
ERC4337_BUNDLER_ETH_CLIENT_URL=http://localhost:8545
ERC4337_BUNDLER_PRIVATE_KEY=c6cbc5ffad570fdad0544d1b5358a36edeb98d163b6567912ac4754e144d4edb
ERC4337_BUNDLER_MAX_BATCH_GAS_LIMIT=12000000
ERC4337_BUNDLER_DEBUG_MODE=true
```

## Running the test suite

Assuming you have your environment properly setup, you can use the following commands to run the QA test suite.

```bash
yarn run test
```
