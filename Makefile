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
	go get github.com/mitchellh/gox

lint:
	golint ./... | grep -v "be unexported"

.PHONY: build
build:
	go build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	ginkgo -r ./...

.PHONY: coverage
coverage:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: docker-build
docker-build:
	docker build -t golang-gluedd .
	docker run -e GOPATH=/ \
	-v $(ROOT_DIR):/src/github.com/mudler/gluedd-cli \
	--workdir /src/github.com/mudler/gluedd-cli \
	--rm -ti golang-gluedd make build

.PHONY: docker-build-arm
docker-build:
	docker build -t golang-gluedd .
	docker run -e GOPATH=/ \
	-e CGO_ENABLED=1 \
	-e CC=arm-linux-gnueabihf-gcc \
	-e GOOS=linux -e GOARCH=arm \
	-v $(ROOT_DIR):/src/github.com/mudler/gluedd-cli \
	--workdir /src/github.com/mudler/gluedd-cli \
	--rm -ti golang-gluedd make build