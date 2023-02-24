export GOEXPERIMENT=arenas

.PHONY: test
test: ## Run tests
	go test -race ./...

.PHONY: bench
bench: ## Run benchmarks
	go test -benchmem -v -race -bench=.

.PHONY: lint
lint: ## Run linter
	golangci-lint run -v ./...

.PHONY: tidy
tidy: ## Run mod tidy
	@echo "Run mod tidy"
	go mod tidy

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
