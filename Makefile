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

# ###############################################
# Version Derivation
#
# Rules and variable definitions used to derive the current version of the
# source code. This information is also used for deriving the type of release
# to perform if `make push` is invoked.
# ###############################################
VERSION=$(shell [ -d .git ] && git describe --tags --abbrev=0 2> /dev/null | sed 's/^v//')
EXACT_TAG=$(shell [ -d .git ] && git describe --exact-match --tags HEAD 2> /dev/null | sed 's/^v//')
ifeq (,$(VERSION))
    VERSION=dev
endif
NOT_RC=$(shell git tag --points-at HEAD | grep -v -e -rc)

ifeq ($(EXACT_TAG),)
    PUSHTYPE=master
else
    ifeq ($(NOT_RC),)
	PUSHTYPE=release-candidate
    else
	PUSHTYPE=release
    endif
endif

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

dockercheck:
ifeq (,$(DOCKER_PATH))
ifeq (,$(findstring $(DOCKER_REQUIRED_VERSION),$(shell docker version)))
ifeq (,$(BYPASS_DOCKER_CHECK))
	$(error "Docker version $(DOCKER_REQUIRED_VERSION) is required.")
endif
endif
endif

# ###############################################
# Building Docker Image
#
# Builds a docker image for TF Encrypted that can be used to deploy and
# test.
# ###############################################
DOCKER_BUILD=docker build -t dropoutlabs/privacyai:$(1) -f Dockerfile $(2) .
docker: Dockerfile dockercheck
	$(call DOCKER_BUILD,latest,)

.PHONY: docker

# ###############################################
# Releasing Docker Images
#
# Using the docker build infrastructure, this section is responsible for
# authenticating to docker hub and pushing built docker containers up with the
# appropriate tags.
# ###############################################
DOCKER_TAG=docker tag dropoutlabs/privacyai:$(1) docker.pkg.github.com/dropoutlabs/privacyai/privacyai:$(2)
DOCKER_PUSH=docker push docker.pkg.github.com/dropoutlabs/privacyai/privacyai:$(1)

docker-tag: dockercheck
	$(call DOCKER_TAG,latest,$(VERSION))

docker-push-tag: dockercheck
	$(call DOCKER_PUSH,$(VERSION))

docker-push-latest: dockercheck
	$(call DOCKER_PUSH,latest)

.PHONY: docker-login docker-push-lateset docker-push-tag docker-tag

# ###############################################
# Targets for pushing docker images
#
# The following are that are called dependent on the push type of the release.
# They define what actions occur depending no whether this is simply a build of
# master (or a branch), release candidate, or a full release.
# ###############################################

# For all builds on the master branch, build the container
docker-push-master: docker

# For all builds on the master branch, with an rc tag
docker-push-release-candidate: releasecheck docker-push-master docker-login docker-tag docker-push-tag

# For all builds on the master branch with a release tag
docker-push-release: docker-push-release-candidate docker-push-latest

# This command calls the right docker push rule based on the derived push type
docker-push: docker-push-$(PUSHTYPE)

.PHONY: docker-push docker-push-release docker-push-release-candidate docker-push-master