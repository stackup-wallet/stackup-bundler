install-dev:
	go install github.com/cosmtrek/air@latest
	go mod tidy

generate-environment:
	go run ./cmd/genenv

generate-entrypoint-pkg:
	abigen --abi=./abi/entrypoint.json --pkg=entrypoint --out=./pkg/entrypoint/bindings.go

fetch-wallet:
	go run ./cmd/fetchwallet

dev-mempool:
	docker-compose -f docker-compose.yml up redis

dev:
	air
