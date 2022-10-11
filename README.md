![](https://i.imgur.com/Kf3qyVJ.png)

# Getting started

A standalone RPC client and bundler for relaying UserOperations to the EntryPoint.

# Running an instance

See the `Client` documentation at [docs.stackup.sh](https://docs.stackup.sh/docs/category/client).

# Contributing

## Prerequisites

- Go 1.19 or later

## Setup

```bash
# Installs https://github.com/cosmtrek/air for live reloading.
# Runs go mod tidy.
make install-dev

# Generates base .env file.
# All variables in this file are required and should be filled.
# Running this command WILL override current .env file.
make generate-environment

# Parses private key in .env file and prints public key and address.
make fetch-wallet
```

## Run RPC server

```bash
make dev
```

# License

Distributed under the GPL-3.0 License. See [LICENSE](./LICENSE) for more information.

# Contact

Feel free to direct any technical related questions to the `dev-hub` channel in the [Stackup Discord](https://discord.gg/FpXmvKrNed).
