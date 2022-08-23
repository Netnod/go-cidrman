TEST?=./...
GOFMT_FILES?=$$(find . -type f -name '*.go' | grep -v 'vendor/')
# Set PKG_NAME when the real package is in a subdir
#PKG_NAME=packagename

default: build

build:
	go install

test:
	go test $(TEST) -timeout=30s -parallel=4

clean:
	go clean -testcache

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

.PHONY: build test clean fmt
