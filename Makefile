#!/usr/bin/env make

.PHONY: build
build:
	@ mkdir -p bin
	go build -o bin/fieldescription ./cmd/fieldescription/main.go

.PHONY: run
run: build
	./bin/fieldescription ./examples/...

.PHONY: test
test:
	go test -v -race -coverprofile cover.out ./pkg/...
	go tool cover -func cover.out