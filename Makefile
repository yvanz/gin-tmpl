PROJECT=gin-demo

GOCMD=CGO_ENABLED=0 go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=$(PROJECT)
BINARY_FILE=cmd/$(BINARY_NAME)/main.go

VERSION=$(shell git describe --abbrev=0 --tags || echo "TagsOrGitNotFound")
COMMITID=$(shell git rev-parse HEAD || echo "GitNotFound")

default: help

## docs: gen swagger doc
docs:
	swag init -g $(BINARY_FILE)

## run: run the code directly
run:
	@$(GOCMD) run $(BINARY_FILE)

## build: build the application
build: clean
	@if [ ! -d "vendor" ]; then $(GOCMD) mod vendor; fi
	$(GOBUILD) -o bin/$(BINARY_NAME) -ldflags " \
			-extldflags '-static' \
			-X 'main.Version=$(VERSION)' \
			-X 'main.Build=$(COMMITID)'" $(BINARY_FILE)

## test: have a test
test:
	$(GOTEST) -v -run=. ./...

## clean: clean the binary
clean:
	@rm -f bin/$(BINARY_NAME)
	$(GOCLEAN) -x

## install: install the application
install:
	@if [ ! -d "vendor" ]; then $(GOCMD) mod vendor; fi
	$(GOCMD) install -mod vendor -ldflags=" \
		-s -w \
		-X 'main.Version=$(VERSION)' \
		-X 'main.Build=$(COMMITID)'" -v ./...

## help: prints this help message
help:
	@echo
	@echo " Choose a command run in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: default run build test clean install help