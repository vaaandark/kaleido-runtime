GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
OUT_DIR = $(PWD)/_output
BIN_DIR = $(OUT_DIR)/bin
BINARY ?= runm

# Git information
GIT_VERSION ?= $(shell git describe --abbrev=10  --match 'v[0-9]*' --dirty='.dirty' --always --tag)
GIT_COMMIT_HASH ?= $(shell git rev-parse HEAD)

.PHONY: help
help:  ## Show this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
all: fmt build ## Build and push the image.

.PHONY: clean
clean:  ## Clean up.
	rm -rf $(OUT_DIR)

.PYONY: fmt
fmt:  ## Format code.
	go fmt ./...

.PYONY: vet
vet:  ## Run go vet.
	@find . -type f -name '*.go'| grep -v "/vendor/" | xargs gofmt -w -s

.PYONY: tidy
tidy:  ## Run go mod tidy.
	@go mod tidy

.PHONY: lint
lint:  ## Run golangci-lint.
	@golangci-lint run --timeout 30m

.PHONY: test
test:  ## Run tests.
	@go test -v ./...

build:  ## Build the binary.
	@echo "building $(BINARY) binary..."
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=$(GOOS) go build -ldflags=$(LDFLAGS) -o $(BIN_DIR)/$(BINARY) cmd/runm/main.go
