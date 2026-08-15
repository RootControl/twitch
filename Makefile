# Binary name. Deliberately not "twitch": that collides with the Twitch CLI
# this tool shells out to, and would shadow it on PATH.
BINARY ?= ttv

# Install location. Defaults to GOBIN, else $(GOPATH)/bin, which the Go
# toolchain already expects to be on your PATH. Override for a system-wide
# install: make install PREFIX=/usr/local/bin (needs sudo).
GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif
PREFIX ?= $(GOBIN)

GOFILES := $(shell find . -name '*.go' -not -path './vendor/*')

.DEFAULT_GOAL := build
.PHONY: build install uninstall run test race vet fmt fmt-check check clean help

## build: compile the binary into the current directory
build: $(BINARY)

$(BINARY): $(GOFILES) go.mod go.sum
	go build -o $(BINARY) .

## install: build and copy the binary to $(PREFIX)
install: build
	@mkdir -p "$(PREFIX)"
	install -m 0755 $(BINARY) "$(PREFIX)/$(BINARY)"
	@echo "installed $(PREFIX)/$(BINARY)"
	@case ":$$PATH:" in \
		*":$(PREFIX):"*) ;; \
		*) echo; echo "warning: $(PREFIX) is not on your PATH."; \
		   echo "add this to your shell profile:"; \
		   echo "    export PATH=\"$(PREFIX):\$$PATH\"" ;; \
	esac

## uninstall: remove the installed binary
uninstall:
	rm -f "$(PREFIX)/$(BINARY)"
	@echo "removed $(PREFIX)/$(BINARY)"

## run: build and run, e.g. make run ARGS="streams -g Rust"
run: build
	@./$(BINARY) $(ARGS)

## test: run all tests
test:
	go test ./...

## race: run all tests under the race detector
race:
	go test -race ./...

## vet: run static analysis
vet:
	go vet ./...

## fmt: format all Go files
fmt:
	gofmt -w $(GOFILES)

## fmt-check: fail if any file needs formatting
fmt-check:
	@unformatted=$$(gofmt -l $(GOFILES)); \
	if [ -n "$$unformatted" ]; then \
		echo "these files need gofmt:"; echo "$$unformatted"; exit 1; \
	fi

## check: fmt-check, vet, and test — run this before pushing
check: fmt-check vet test

## clean: remove build artifacts
clean:
	rm -f $(BINARY)
	go clean

## help: list available targets
help:
	@echo "usage: make [target]"
	@echo
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | awk -F': ' '{printf "  %-12s %s\n", $$1, $$2}'
	@echo
	@echo "variables:"
	@echo "  BINARY=$(BINARY)"
	@echo "  PREFIX=$(PREFIX)"
