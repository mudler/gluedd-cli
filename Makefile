NAME ?= gluedd-cli
PACKAGE_NAME ?= $(NAME)
PACKAGE_CONFLICT ?= $(PACKAGE_NAME)-beta
REVISION := $(shell git rev-parse --short HEAD || echo unknown)
VERSION := $(shell git describe --tags || echo dev)
VERSION := $(shell echo $(VERSION) | sed -e 's/^v//g')
ITTERATION := $(shell date +%s)

EXTENSIONS ?=
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: test build

help:
	# make all => deps test lint build
	# make deps - install all dependencies
	# make test - run project tests
	# make lint - check project code style
	# make build - build project for all supported OSes
deps:
	go env
	# Installing dependencies...
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega/...

lint:
	golint ./... | grep -v "be unexported"

.PHONY: build
build:
	go build -ldflags "-extldflags=-Wl,--allow-multiple-definition"


.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	ginkgo -r ./...

.PHONY: coverage
coverage:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

