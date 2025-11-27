LC_ALL=C
export LC_ALL

.DEFAULT_GOAL := all

.PHONY: all
all: build

# set the shell to bash in case some environments use sh
SHELL := /usr/bin/env bash

# include the common make file
SELF_DIR := $(shell pwd)

# Can be used or additional go build flags
BUILDFLAGS ?=
LDFLAGS ?=
TAGS ?=

# Set GOBIN
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# the root go project
GO_PROJECT=github.com/baptistegh/go-lakekeeper

# CGO_ENABLED value
CGO_ENABLED_VALUE=0

COMMIT_SHA := $(shell git rev-parse HEAD)
DATE := $(shell git log -1 --format=%cI)

ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --dirty --always --tags | sed 's/-/./2' | sed 's/-/./2' )
endif
export VERSION

GOHOSTOS=linux
GOHOSTARCH := $(shell go env GOHOSTARCH)
HOST_PLATFORM := $(GOHOSTOS)_$(GOHOSTARCH)

# REAL_HOST_PLATFORM is used to determine the correct url to download the various binary tools from and it does not use
# HOST_PLATFORM which is used to build the program.
REAL_HOST_PLATFORM=$(shell go env GOHOSTOS)_$(GOHOSTARCH)

GO := go
GOHOST := GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) $(GO)
GO_VERSION := $(shell $(GO) version | sed -ne 's/[^0-9]*\(\([0-9]\.\)\{0,4\}[0-9][^.]\).*/\1/p')
GO_FULL_VERSION := $(shell $(GO) version)

BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_FULL_COMMIT := $(shell git rev-parse HEAD)
GIT_TREE_STATE := $(shell if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi)

GO_BUILDFLAGS=$(BUILDFLAGS)
GO_LDFLAGS=-X github.com/baptistegh/go-lakekeeper/pkg/version.buildDate=$(NOW) -X github.com/baptistegh/go-lakekeeper/pkg/version.gitCommit=$(FULL_COMMIT) -X github.com/baptistegh/go-lakekeeper/pkg/version.gitTreeState=$(GIT_TREE_STATE) $(LDFLAGS)
GO_TAGS=$(TAGS)

GO_COMMON_FLAGS = $(GO_BUILDFLAGS) -tags '$(GO_TAGS)' -ldflags '$(GO_LDFLAGS)'

GO_PACKAGES := ./pkg/...

BIN_DIR := $(SELF_DIR)/bin
DIST_DIR := $(SELF_DIR)/dist

CONTAINER_ENGINE ?= docker
CONTAINER_COMPOSE_ENGINE ?= $(shell $(CONTAINER_ENGINE) compose version >/dev/null 2>&1 && echo '$(CONTAINER_ENGINE) compose' || echo '$(CONTAINER_ENGINE)-compose')
DOCKER ?= docker

ENV_FILE := $(SELF_DIR)/.env

$(BIN_DIR):
	@echo === creating $(BIN_DIR)
	@mkdir -p $(BIN_DIR)

YQ_VERSION := v4.45.1
YQ := $(BIN_DIR)/yq-$(YQ_VERSION)
$(YQ): | $(BIN_DIR)
	@echo === installing yq $(YQ_VERSION) $(REAL_HOST_PLATFORM)
	@curl -s -JL https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(REAL_HOST_PLATFORM) -o $(YQ)
	@chmod +x $(YQ)
	@echo Installed yq version $(YQ_VERSION) in $(YQ)

.PHONY: build
build: build.common
	@echo === go build dist/lkctl
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) build -a $(GO_COMMON_FLAGS) -o $(DIST_DIR)/lkctl ./cmd

.PHONY: build.common
build.common: $(YQ) mod fmt vet test

.PHONY: mod
mod: 
	@$(GO) mod tidy

.PHONY: test
test: ## Runs unit tests.
	@echo === go test unit-tests
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) test -v -cover -coverprofile=coverage.txt $(GO_COMMON_FLAGS) $(GO_PACKAGES)

LAKEKEEPER_VERSION ?= latest-main
.PHONY: test-integration
test-integration: $(ENV_FILE) ## Runs integration tests.
	@echo === ./run-tests.sh
	CONTAINER_COMPOSE_ENGINE="$(CONTAINER_COMPOSE_ENGINE)" LAKEKEEPER_VERSION="$(LAKEKEEPER_VERSION)" ./run-tests.sh

GORELEASER := $(BIN_DIR)/goreleaser
$(GORELEASER): | $(BIN_DIR)
	@echo === installing goreleaser
	@GOBIN=$(BIN_DIR) $(GO) install github.com/goreleaser/goreleaser/v2@latest

.PHONY: snapshot
snapshot: $(GORELEASER)
	@echo === goreleaser snapshot
	@GIT_TREE_STATE=$(GIT_TREE_STATE) $(GORELEASER) --clean --snapshot --skip sign

$(ENV_FILE):
	@echo === creating integration tests environments
	@echo 'LAKEKEEPER_BASE_URL="http://localhost:8181"' > $(ENV_FILE)
	@echo 'LAKEKEEPER_TOKEN_URL="http://localhost:30080/realms/iceberg/protocol/openid-connect/token"' >> $(ENV_FILE)
	@echo 'LAKEKEEPER_SCOPE="lakekeeper"' >> $(ENV_FILE)
	@echo 'LAKEKEEPER_CLIENT_ID="lakekeeper-admin"' >> $(ENV_FILE)
	@echo 'LAKEKEEPER_CLIENT_SECRET="KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL"' >> $(ENV_FILE)

.PHONY: vet
vet:
	@echo === go vet
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) vet $(GO_COMMON_FLAGS) ./...

.PHONY: fmt
fmt:
	@echo === golangci-lint fix
	@$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --fix ./...

.PHONY: lint
lint:
	@echo === golangci-lint
	@$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run ./...

.PHONY: validate
validate: vet lint

.PHONY: clean
clean:
	@rm -fr $(BIN_DIR)
	@rm -fr coverage.txt
	@rm -fr $(ENV_FILE)
	@$(CONTAINER_COMPOSE_ENGINE) down --volumes