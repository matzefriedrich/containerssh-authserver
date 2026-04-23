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

.PHONY: all
all: build

.PHONY: build
build: generate
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(CMD_PATH)

.PHONY: test
test: generate
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
	go install github.com/matzefriedrich/parsley/cmd/parsley-cli@${PARSLEY_CLI_VERSION}

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
