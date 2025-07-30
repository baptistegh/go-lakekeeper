# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Linux doesn't guarantee file ordering, so sort the files to make sure order is deterministic.
# And in order to handle file paths with spaces, it's easiest to read the file names into an array.
# Set locale `LC_ALL=C` because different OSes have different sort behavior;
# `C` sorting order is based on the byte values,
# Reference: https://blog.zhimingwang.org/macos-lc_collate-hunt
LC_ALL=C
export LC_ALL

.PHONY: all
all: build
.DEFAULT_GOAL := all

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

GO_BUILDFLAGS=$(BUILDFLAGS)
GO_LDFLAGS=$(LDFLAGS)
GO_TAGS=$(TAGS)

GO_COMMON_FLAGS = $(GO_BUILDFLAGS) -tags '$(GO_TAGS)' -ldflags '$(GO_LDFLAGS)'

GO_PACKAGES := ./pkg/...

BIN_DIR := $(SELF_DIR)/bin
DIST_DIR := $(SELF_DIR)/dist

CONTAINER_COMPOSE_ENGINE ?= $(shell docker compose version >/dev/null 2>&1 && echo 'docker compose' || echo 'docker-compose')

ENV_FILE := $(SELF_DIR)/.env

$(BIN_DIR):
	@echo === creating $(BIN_DIR)
	@mkdir -p $(BIN_DIR)

YQ_VERSION := v4.45.1
YQ := $(BIN_DIR)/yq-$(YQ_VERSION)
$(YQ): $(BIN_DIR)
	@echo === installing yq $(YQ_VERSION) $(REAL_HOST_PLATFORM)
	@curl -s -JL https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(REAL_HOST_PLATFORM) -o $(YQ)
	@chmod +x $(YQ)
	@echo Installed yq version $(YQ_VERSION) in $(YQ)

GOLANGCI_LINT_VERSION ?= $(strip $(shell $(YQ) .jobs.lint.steps[2].with.version .github/workflows/test.yml))
GOLANGCI_LINT ?= $(BIN_DIR)/golangci-lint
$(GOLANGCI_LINT): $(BIN_DIR)
	@echo === installing golangci-lint
	@mkdir -p $(BIN_DIR)/tmp
	@curl -sL https://github.com/golangci/golangci-lint/releases/download/$(GOLANGCI_LINT_VERSION)/golangci-lint-$(patsubst v%,%,$(GOLANGCI_LINT_VERSION))-$(shell go env GOHOSTOS)-$(GOHOSTARCH).tar.gz | tar -xz -C $(BIN_DIR)/tmp
	@mv $(BIN_DIR)/tmp/golangci-lint-$(patsubst v%,%,$(GOLANGCI_LINT_VERSION))-$(shell go env GOHOSTOS)-$(GOHOSTARCH)/golangci-lint $(GOLANGCI_LINT)
	@rm -fr $(BIN_DIR)/tmp

ADD_LICENSE := $(BIN_DIR)/addlicense
$(ADD_LICENSE): $(BIN_DIR)
	@echo === installing addlicense
	@GOBIN=$(BIN_DIR) $(GO) install github.com/google/addlicense@latest

.PHONY: build
build: build.common
	@echo === go build $(DIST_DIR)/lkctl
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) build -a $(GO_COMMON_FLAGS) -o $(DIST_DIR)/lkctl main.go

.PHONY: build.common
build.common: $(YQ)
	@$(GOHOST) mod tidy
	@$(MAKE) fmt
	@$(MAKE) validate
	@$(MAKE) test

.PHONY: test
test: ## Runs unit tests.
	@echo === go test unit-tests
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) test -v -cover -coverprofile=coverage.txt $(GO_COMMON_FLAGS) $(GO_PACKAGES)

LAKEKEEPER_VERSION ?= latest-main
.PHONY: test-integration
test-integration: $(ENV_FILE) ## Runs integration tests.
	@echo === ./run-tests.sh
	CONTAINER_COMPOSE_ENGINE="$(CONTAINER_COMPOSE_ENGINE)" LAKEKEEPER_VERSION="$(LAKEKEEPER_VERSION)" ./run-tests.sh

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
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GOHOST) vet $(GO_COMMON_FLAGS) ./...

.PHONY: fmt
fmt: license $(GOLANGCI_LINT)
	@echo === go fmt fix
	@$(GOLANGCI_LINT) run --fix ./...

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo === go fmt
	@$(GOLANGCI_LINT) run ./...

.PHONY: validate
validate: license-check vet lint

COPYRIGHT ?= Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
.PHONY: license
license: $(ADD_LICENSE)
	@echo ===  addlicense
	@$(ADD_LICENSE) -l apache -c "$(COPYRIGHT)" .

.PHONY: license-check
license-check: $(ADD_LICENSE)
	@echo === addlicense-check
	@$(ADD_LICENSE) -check -l apache -c "$(COPYRIGHT)" .

.PHONY: clean
clean:
	@rm -fr $(BIN_DIR)
	@rm -fr coverage.txt
	@rm -fr $(ENV_FILE)
	@$(CONTAINER_COMPOSE_ENGINE) down --volumes