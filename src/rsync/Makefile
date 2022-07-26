# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2020 Intel Corporation

export GO111MODULE=on

all: clean
	GOOS=linux GOARCH=amd64 \
	go build -tags static,system_libgit2,rsync -o ../../bin/rsync/rsync ./cmd/main.go

# The following is done this way as each patch on CI runs build and each merge runs deploy. So for build we don't need to build binary and hence
# no need to create a static binary with additional flags. However, for generating binary, additional build flags are necessary. This if used with
# mock plugin errors out for unit tests. So the seperation avoids the error.

build: clean test cover
deploy: build

.PHONY: test
test: clean
	@go test -race ./...

format:
	@go fmt ./...

clean:
	@rm -f ../../bin/rsync/rsync coverage.html coverage.out

.PHONY: cover
cover:
	@go test -race ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
