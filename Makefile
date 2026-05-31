.PHONY: run build test dev

run:
	go run ./cmd/main.go

build:
	go build -o bin/BBingyan ./cmd/main.go

test:
	go test ./internal/... -v -cover

dev:
	air -c .air.toml
