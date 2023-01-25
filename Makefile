NAME := nsxctl
RELEASE_DIR := build
BUILD_TARGETS := build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-arm64
GOVERSION = $(shell go version)
THIS_GOOS = $(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH = $(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS = $(THIS_GOOS)
GOARCH = $(THIS_GOARCH)
VERSION = $(patsubst "%",%,$(lastword $(shell grep 'const version' main.go)))
REVISION = $(shell git rev-parse HEAD)

.PHONY: fmt build clean

fmt: ## format
	go fmt

lint: ## Examine source code and lint
	go vet ./...
	golint -set_exit_status ./...

all: $(BUILD_TARGETS) ## build for all platform

build: $(RELEASE_DIR)/nsxctl_$(GOOS)_$(GOARCH) ## build nsxctl

build-linux-amd64: ## build AMD64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm64: ## build ARM64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-darwin-amd64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-arm64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=arm64

build-windows-amd64: ## build AMD64 windows binary
	@$(MAKE) build GOOS=windows GOARCH=amd64

build-windows-arm64: ## build AMD64 windows binary
	@$(MAKE) build GOOS=windows GOARCH=arm64

$(RELEASE_DIR)/nsxctl_$(GOOS)_$(GOARCH): ## Build nsx command-line client
	@printf "\e[32m"
	@echo "==> Build nsxctl for ${GOOS}-${GOARCH}"
	@printf "\e[90m"
	@GO111MODULE=on go build -tags netgo -ldflags "-X github.com/hichtakk/nsxctl/cmd.revision=${REVISION}" -a -v -o $(RELEASE_DIR)/nsxctl_$(GOOS)_$(GOARCH) ./main.go
	@printf "\e[m"

clean: ## Clean up built files
	@printf "\e[32m"
	@echo '==> Remove built files ./build/...'
	@printf "\e[90m"
	@ls -1 ./build
	@rm -rf build/*
	@printf "\e[m"

rebuild: clean build