.PHONY: test

install-deps:
	go mod download

test:
	go test ./...