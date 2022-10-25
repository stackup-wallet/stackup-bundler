install-dev:
	go install github.com/cosmtrek/air@latest
	go mod tidy

generate-environment:
	go run ./scripts/genenv

generate-entrypoint-pkg:
	abigen --abi=./abi/entrypoint.json --pkg=entrypoint --out=./pkg/entrypoint/bindings.go

fetch-wallet:
	go run ./scripts/fetchwallet

dev-mempool:
	docker-compose -f docker-compose.yml up redis

dev-client:
	air -c .air.client.toml
