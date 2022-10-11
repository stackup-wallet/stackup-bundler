install-dev:
	go install github.com/cosmtrek/air@latest
	go mod tidy

generate-environment:
	go run ./cmd/genenv

fetch-wallet:
	go run ./cmd/fetchwallet

dev:
	air
