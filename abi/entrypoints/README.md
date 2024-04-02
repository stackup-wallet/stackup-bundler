# HOW-TO

## Prep

```shell
go get -u github.com/ethereum/go-ethereum
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

> Note: `abigen` is a tool that is part of the go-ethereum project. It is used to generate Go bindings for Ethereum contracts.

### test

```shell
abigen --version
```

## Generate Go bindings

```shell
abigen --abi=abi/entrypoints/v0.6.json --pkg=entrypoint --out=./pkg/entrypoint/bindings/v06/bindings.go
abigen --abi=abi/entrypoints/v0.7.json --pkg=entrypoint --out=./pkg/entrypoint/bindings/v07/bindings.go
```

