# Makefile for Client-Server Golang Project

.PHONY: all build client server fmt vet clean test

# Output binaries
CLIENT_BIN := bin/client
SERVER_BIN := bin/server

# Point directly to package paths
CLIENT_PKG := ./client
SERVER_PKG := ./server

# Version info compiled into the binaries. GIT_REF falls back to 'unknown'
# outside a git working copy (e.g. an unpacked source tarball).
APP_VERSION := $(shell . ./common.sh && echo $$APP_VERSION)
GIT_REF := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -X rpm-signing/common.AppVersion=$(APP_VERSION) -X rpm-signing/common.GitRef=$(GIT_REF)

# Default target runs code quality checks and builds both binaries
all: fmt vet build test

all-stripped: fmt vet build test
    LDFLAGS += -s -w

# Build both client and server
build: client server

build-stripped: client server
    LDFLAGS += -s -w

# Build client binary
client:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o $(CLIENT_BIN) $(CLIENT_PKG)

# Build server binary
server:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o $(SERVER_BIN) $(SERVER_PKG)

# Run go fmt across the project
fmt:
	go fmt ./...

# Run go vet across the project
vet:
	go vet ./...

test:
	go test ./...

# Remove built binaries
clean:
	rm -rf bin/