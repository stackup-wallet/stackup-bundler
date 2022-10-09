setup-dev:
	go install github.com/cosmtrek/air@latest
	go mod tidy

dev-rpc:
	air -c .air-rpc.toml
