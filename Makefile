GO=go

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

packages := $(shell go list ./...|grep -v /vendor/)

.PHONY: check test lint

test: check
	@$(GO) test -race $(packages) -v -coverprofile=.coverage.out
	@$(GO) tool cover -func=.coverage.out
	@rm -f .coverage.out

check:
	@$(GO) vet -composites=false $(packages)

lint:
	@golangci-lint run ./...

doc:
	@godoc -http=localhost:8098
