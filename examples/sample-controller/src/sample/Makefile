# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

export GO111MODULE=on

all: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -tags netgo -o ../../bin/sample/sample ./cmd/main.go
	# Please modify the output directory accordingly.

build: clean test cover

deploy: build

.PHONY: test
test: clean
	@go test -race ./...

format:
	@go fmt ./...

clean:
	@rm -f ../../bin/sample/sample coverage.html coverage.out
	# Please modify accordingly.

.PHONY: cover
cover:
	@go test -race ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
