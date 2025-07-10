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

# packages that will be compiled into binaries
CLIENT_PACKAGES = $(GO_PROJECT)/cmd/rook

# the root go project
GO_PROJECT=github.com/baptistegh/go-lakekeeper

# CGO_ENABLED value
CGO_ENABLED_VALUE=0

COMMIT_SHA := $(shell git rev-parse HEAD)

DATE := $(shell git log -1 --format=%cI)

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --dirty --always --tags | sed 's/-/./2' | sed 's/-/./2' )
endif
export VERSION

# inject the version and details the version package using the -X linker flag
LDFLAGS += -X $(GO_PROJECT)/pkg/version.Version=$(VERSION) \
           -X $(GO_PROJECT)/pkg/version.Commit=$(COMMIT_SHA) \
           -X $(GO_PROJECT)/pkg/version.Date=$(DATE)

GOHOSTOS=linux
GOHOSTARCH := $(shell go env GOHOSTARCH)
HOST_PLATFORM := $(GOHOSTOS)_$(GOHOSTARCH)

# REAL_HOST_PLATFORM is used to determine the correct url to download the various binary tools from and it does not use
# HOST_PLATFORM which is used to build the program.
REAL_HOST_PLATFORM=$(shell go env GOHOSTOS)_$(GOHOSTARCH)

GO := go
GOHOST := GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) go
GO_VERSION := $(shell $(GO) version | sed -ne 's/[^0-9]*\(\([0-9]\.\)\{0,4\}[0-9][^.]\).*/\1/p')
GO_FULL_VERSION := $(shell $(GO) version)

GO_BUILDFLAGS=$(BUILDFLAGS)
GO_LDFLAGS=$(LDFLAGS)
GO_TAGS=$(TAGS)

GO_COMMON_FLAGS = $(GO_BUILDFLAGS) -tags '$(GO_TAGS)' -ldflags '$(GO_LDFLAGS)'

GO_PACKAGES := ./cmd/... ./pkg/...

GO_BUILD_TARGET := ./cmd/lakekeeper/main.go

BIN_DIR := $(SELF_DIR)/bin

GO_TEST_OUTPUT := $(SELF_DIR)/.output

build: build.common test ## Only build for linux platform

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(GO_TEST_OUTPUT):
	@mkdir -p $(GO_TEST_OUTPUT)

YQ_VERSION := v4.45.1
YQ := $(BIN_DIR)/yq-$(YQ_VERSION)
$(YQ): $(BIN_DIR)
	@echo === installing yq $(YQ_VERSION) $(REAL_HOST_PLATFORM)
	@mkdir -p $(BIN_DIR)
	@curl -s -JL https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(REAL_HOST_PLATFORM) -o $(YQ)
	@chmod +x $(YQ)

GOLANGCI_LINT_VERSION ?= $(strip $(shell $(YQ) .jobs.golangci.steps[2].with.version .github/workflows/lint.yml))
GOLANGCI_LINT ?= $(BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)
$(GOLANGCI_LINT):
	@echo === installing golangci-lint-$(GOLANGCI_LINT_VERSION)
	@mkdir -p $(BIN_DIR)/tmp
	@curl -sL https://github.com/golangci/golangci-lint/releases/download/$(GOLANGCI_LINT_VERSION)/golangci-lint-$(patsubst v%,%,$(GOLANGCI_LINT_VERSION))-$(shell go env GOHOSTOS)-$(GOHOSTARCH).tar.gz | tar -xz -C $(BIN_DIR)/tmp
	@mv $(BIN_DIR)/tmp/golangci-lint-$(patsubst v%,%,$(GOLANGCI_LINT_VERSION))-$(shell go env GOHOSTOS)-$(GOHOSTARCH)/golangci-lint $(GOLANGCI_LINT)
	@rm -fr $(BIN_DIR)/tmp

build.common: $(YQ) $(BIN_DIR)
	@$(GO) mod tidy
	@$(MAKE) validate
	@CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) build $(GO_COMMON_FLAGS) -o $(BIN_DIR)/lakekeeper $(GO_BUILD_TARGET)

.PHONY: test
test: $(GO_TEST_OUTPUT) ## Runs unit tests.
	@echo === go test unit-tests
	@mkdir -p $(GO_TEST_OUTPUT)
	CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GO) test -v -cover -coverprofile=coverage.txt $(GO_COMMON_FLAGS) $(GO_PACKAGES)

.PHONY: vet
vet:
	@echo === go vet
	CGO_ENABLED=$(CGO_ENABLED_VALUE) $(GOHOST) vet $(GO_COMMON_FLAGS) ./...

.PHONY: validate
validate: vet lint

.PHONY: fmt
fmt: $(GOLANGCI_LINT)
	@echo === go fmt fix
	@$(GOLANGCI_LINT) run --fix

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo === go fmt
	@$(GOLANGCI_LINT) run
