# Global Variables used across many different rule types
PREFIX?=
GOOS_OVERRIDE?=
DOCKER_REQUIRED_VERSION=18.

SRC=$(shell find . -type f -name '*.go')
rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) \
	$(filter $(subst *,%,$2),$d))

all: build

.PHONY: all

GO_REQUIRED_VERSION=1.13.
GOLINT_REQUIRED_VERSION=v1.22.2

VERSION=0.1.0

gocheck:
ifeq (,$(findstring $(GO_REQUIRED_VERSION),$(shell go version)))
ifeq (,$(BYPASS_GO_CHECK))
	$(error "Go Version $(GO_REQUIRED_VERSION) is required.")
endif
endif

golangci-lint-check:
ifeq (,$(findstring $(GOLINT_REQUIRED_VERSION),$(shell golangci-lint --version)))
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLINT_REQUIRED_VERSION)
endif

bootstrap:
	go mod download -json

.PHONY: bootstrap gocheck

lint: golangci-lint-check
	golangci-lint run

GO_BUILD=go build -v -i

$(PREFIX)bin/privacy: $(call rwildcard,,*.go)
	$(GOOS_OVERRIDE) $(GO_BUILD) -o $@

build: $(PREFIX)bin/privacy

test:
	go test -v $$(go list ./...)

fmt:
	gofmt -s -l -w $(SRC)

clean:
	go clean
	rm $(PREFIX)bin/privacy

ci: test lint

.PHONY: lint build fmt

docker:
	docker build . --file Dockerfile --tag dropoutlabs/privacyai:$(shell date +%s)