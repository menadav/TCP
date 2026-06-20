GO       ?= go
BIN_DIR  := bin
SERVER   := $(BIN_DIR)/tap-server
CLI      := $(BIN_DIR)/tap-client
GUI      := $(BIN_DIR)/tap-client-gui

.PHONY: all install build run-server run-client run-client-gui lint fmt fmt-fix vet clean

all: build

install:
	$(GO) mod download

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(SERVER) ./cmd/server
	$(GO) build -o $(CLI) ./cmd/client-cli
	$(GO) build -o $(GUI) ./cmd/client-gui

run-server:
	$(GO) run ./cmd/server

run-client:
	$(GO) run ./cmd/client-cli

run-client-gui:
	$(GO) run ./cmd/client-gui

lint: fmt vet

fmt:
	@unformatted="$$(gofmt -l src cmd)"; \
	if [ -n "$$unformatted" ]; then \
		echo "Files not gofmt-clean:"; echo "$$unformatted"; \
		exit 1; \
	fi

fmt-fix:
	gofmt -w src cmd

vet:
	$(GO) vet ./...

clean:
	$(GO) clean
	rm -rf $(BIN_DIR)
