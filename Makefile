SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GO ?= go
TEST := $(GO) test
TEST_FLAGS ?= -v
TEST_TARGET ?= ./...
GO111MODULE = on
PROJECT_NAME := $(shell basename $(PWD))

.PHONY: test coverage clean download

download: go.sum

go.sum: go.mod
	$(GO) mod tidy

test: go.sum clean
	@
	$(TEST) $(TEST_FLAGS) -cover $(TEST_TARGET) -json > tests.json || true
	$(GO) tool tparse -file tests.json

coverage: go.sum clean
	@mkdir ./_coverage
	$(TEST) $(TEST_FLAGS) -covermode=count -coverpkg github.com/go-ap/filters,github.com/go-ap/filters/index -args -test.gocoverdir="$(PWD)/_coverage" . #> /dev/null || true
	pushd index
	$(TEST) $(TEST_FLAGS) -covermode=count -coverpkg github.com/go-ap/filters,github.com/go-ap/filters/index -args -test.gocoverdir="$(PWD)/_coverage" . #> /dev/null || true
	popd
	$(GO) tool covdata percent -i=./_coverage/ -o $(PROJECT_NAME).coverprofile
	@$(RM) -r ./_coverage

clean:
	@$(RM) -r ./_coverage
	@$(RM) -v *.coverprofile
	@$(RM) -v tests.json

