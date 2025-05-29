.PHONY: build test test-verbose

build:
	go build -o build/runes

test:
	go test ./...

test-verbose:
	go test -v ./...