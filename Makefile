.PHONY: all build test lint lint-fix generate install-tools clean

PACKAGE_NAME=github.com/matzefriedrich/containerssh-authserver
CI_COMMIT_SHORT_SHA?=$(shell git rev-parse --short=8 HEAD 2>/dev/null || echo "unknown")
APP_RELEASE_DATE?=$(shell date -u +%Y-%m-%d)
APP_RELEASE?=matzefriedrich/containerssh-authserver
APP_VERSION?=0.1.0

LDFLAGS=-X $(PACKAGE_NAME)/internal.CommitSha=$(CI_COMMIT_SHORT_SHA) \
        -X $(PACKAGE_NAME)/internal.Version=$(APP_VERSION) \
        -X $(PACKAGE_NAME)/internal.ReleaseDate=$(APP_RELEASE_DATE) \
        -X $(PACKAGE_NAME)/internal.ReleaseName=$(APP_RELEASE)

BINARY_NAME=authserver
CMD_PATH=./cmd/authserver

GOLANGCI_LINT_VERSION=v2.11.4
PARSLEY_CLI_VERSION=v1.4.0

all: install-tools generate lint build

build: generate
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(CMD_PATH)

test: generate
	go test ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

generate:
	go generate ./...

install-tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
	go install github.com/matzefriedrich/parsley/cmd/parsley-cli@${PARSLEY_CLI_VERSION}

clean:
	rm -f $(BINARY_NAME)
