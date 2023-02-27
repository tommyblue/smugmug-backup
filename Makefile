# This Makefile can be used to obtain a Linux binary through docker (eg: from OSX)
# For normal development, Docker is not required. See README.md for build instructions.
.PHONY: build help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build smugmug-backup binary in the local env
	GOFLAGS=-mod=vendor go build -v ./cmd/smugmug-backup/

test-short: ## Run tests with -short flag in the local env
	GOFLAGS=-mod=vendor gotestsum --format testname -- -short -shuffle=on -cover -race -v -count=1 ./...

test: ## Run tests in the local env
	GOFLAGS=-mod=vendor gotestsum --format testname -- -shuffle=on -cover -race -v -count=1 ./...

gofmt: ## Run gofmt locally without overwriting any file
	gofmt -d -s $$(find . -name '*.go' | grep -v vendor)

gofmt-write: ## Run gofmt locally overwriting files
	gofmt -w -s $$(find . -name '*.go' | grep -v vendor)

govet: ## Run go vet on the project
	GOFLAGS=-mod=vendor go vet ./...

docker: ## Build docker image with current version of the code
	docker build -t smugmug .

docker-test-short: docker ## Run tests with -short flag with Docker
	docker run -t -v "$$PWD:/go/src/smugmug-backup" -e 'GOFLAGS=-mod=vendor' smugmug go test -short -race ./...

docker-test: docker ## Run tests with Docker
	docker run -t -v "$$PWD:/go/src/smugmug-backup" -e 'GOFLAGS=-mod=vendor' smugmug go test -race ./...

docker-build: docker ## Build the smugmug-backup linux binary /smugmug-backup-linux using Docker
	docker run -t -v "$$PWD:/go/src/smugmug-backup" -e 'GOFLAGS=-mod=vendor' smugmug go build -v -o smugmug-backup-linux ./cmd/smugmug-backup/
