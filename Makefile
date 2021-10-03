NAME := nsxctl
RELEASE_DIR := build

.PHONY: fmt build clean

fmt: ## format
	go fmt

build: ## build nsxctl
	go build -a -v -o ${RELEASE_DIR}/nsxctl ./main.go

clean: ## clean up files
	rm -rf ./build