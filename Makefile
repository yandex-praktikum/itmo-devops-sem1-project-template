# TOOLS_PATH defines path to Golang-based utility binaries.
TOOLS_PATH=bin/tools

# Delegate calling Golang utilities to locally installed stuff in $TOOLS_PATH.
gci=${TOOLS_PATH}/gci
gofumpt=${TOOLS_PATH}/gofumpt
golangci-lint=${TOOLS_PATH}/golangci-lint

.PHONY: FORCE
.DEFAULT_GOAL := build

$(gci): Makefile
	GOBIN=`pwd`/$(TOOLS_PATH) go install github.com/daixiang0/gci@v0.13.0

$(gofumpt): Makefile
	GOBIN=`pwd`/$(TOOLS_PATH) go install mvdan.cc/gofumpt@v0.6.0

$(golangci-lint): Makefile
	GOBIN=`pwd`/$(TOOLS_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

setup: FORCE $(gci) $(gofumpt) $(golangci-lint)

fmt: FORCE fmt/sources fmt/imports ## Format all the stuff

fmt/sources: FORCE $(gofumpt) ## Format the source files
	$(gofumpt) -l -w .

fmt/imports: FORCE $(gci) ## Format imports using gci tool
	$(gci) list . | xargs $(gci) write --skip-generated -s standard -s 'prefix(project_sem)' -s default --custom-order

lint: FORCE lint/sources ## Run all linters

lint/sources: FORCE $(golangci-lint) ## Lint the source files
	$(golangci-lint) run --timeout 6m