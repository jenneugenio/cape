PREFIX?=
GOOS_OVERRIDE?=
SRC=$(shell find . -type f -name '*.go')

all: ci

.PHONY: all

# ###############################################
# Required Binaries & Versioning Checks
#
# The following section lists out all of the system dependencies (and their
# versions) that this Makefile depends upon on-top of standard tooling such as
# Git.
# ###############################################
DOCKER_REQUIRED_VERSION=18.
GO_REQUIRED_VERSION=1.13.
GOLINT_REQUIRED_VERSION=v1.22.2
HELM_REQUIRED_VERSION=v3.0.
GSUTIL_REQUIRED_VERSION=4.47

CURRENT_GO_VERSION := $(shell go version 2> /dev/null)
gocheck:
ifndef CURRENT_GO_VERSION
ifeq (,$(findstring $(GO_REQUIRED_VERSION),$(CURRENT_GO_VERSION)))
ifeq (,$(BYPASS_GO_CHECK))
	$(error "Go Version $(GO_REQUIRED_VERSION) is required, found $(CURRENT_GO_VERSION)")
endif
endif
endif

CURRENT_DOCKER_VERSION := $(shell docker version 2> /dev/null)
dockercheck:
ifeq (,$(DOCKER_PATH))
ifeq (,$(shell command -v docker 2> /dev/null))
ifeq (,$(findstring $(DOCKER_REQUIRED_VERSION),$(CURRENT_DOCKER_VERSION)))
ifeq (,$(BYPASS_DOCKER_CHECK))
	$(error "Docker version $(DOCKER_REQUIRED_VERSION) is required, found $(CURRENT_DOCKER_VERSION)")
endif
endif
endif
endif

CURRENT_HELM_VERSION := $(shell helm version 2> /dev/null)
helmcheck:
ifeq (,$(shell command -v helm 2> /dev/null))
ifeq (,$(findstring $(HELM_REQUIRED_VERSION),$(CURRENT_HELM_VERSION)))
ifeq (,$(BYPASS_HELM_CHECK))
	$(error "Helm version $(HELM_REQUIRED_VERSION) is required, found $(CURRENT_HELM_VERSION).")
endif
endif
endif

CURRENT_GSUTIL_VERSION := $(shell gsutil version 2> /dev/null)
gsutilcheck:
ifeq (,$(shell command -v gsutil 2> /dev/null))
ifeq (,$(findstring $(GSUTIL_REQUIRED_VERSION),$(CURRENT_GSUTIL_VERSION)))
ifeq (,$(BYPASS_GSUTIL_CHECK))
	$(error "gsutil version $(GSUTIL_REQUIRED_VERSION) is required, found $(CURRENT_GSUTIL_VERSION)")
endif
endif
endif

.PHONY: gocheck dockercheck

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
    VERSION=0.0.1
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

# ###############################################
# Setup and Teardown
#
# This section of the makefile focuses on the rules needed to setup an
# environment to develop with this codebase and then subsequently start from
# scratch.
# ###############################################

bootstrap: download install-tools
	$(info All dependencies and tooling have been installed!)

download: gocheck
	go mod download

install-tools: gocheck download tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean: gocheck
	go clean
	rm $(PREFIX)bin/privacy

.PHONY: bootstrap clean

bootstrap-local-dev: bootstrap-helm

bootstrap-helm: helm-install helm-add-stable helm-update

helm-install:
ifeq (, $(shell which helm))
	curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
	chmod 700 get_helm.sh
	./get_helm.sh
endif

helm-add-stable: helmcheck
	helm repo add stable https://kubernetes-charts.storage.googleapis.com

helm-update: helmcheck
	helm repo update

.PHONY: bootstrap-local-dev

# ###############################################
# Testing, Building and Formatting
#
# This section of the makefile focuses on the rules needed to test and build
# code contained within this repository.
# ###############################################

lint: gocheck
	golangci-lint run

GO_BUILD=go build -v -i
$(PREFIX)bin/privacy: gocheck $(SRC)
	$(GOOS_OVERRIDE) $(GO_BUILD) -o $@

build: $(PREFIX)bin/privacy

unit: gocheck
	go test -v ./...

CAPE_DB_URL?="postgres://postgres:dev@localhost:5432/postgres?sslmode=disable"
integration: gocheck
	CAPE_DB_URL=$(CAPE_DB_URL) go test -v ./... -tags=integration

test: integration

fmt: gocheck
	gofmt -s -l -w $(SRC)

ci: lint build test docker

.PHONY: lint build fmt test ci

# ###############################################
# Building Docker Image
#
# Builds a docker image for TF Encrypted that can be used to deploy and
# test.
# ###############################################
DOCKER_BUILD=docker build -t dropoutlabs/$(1):$(2) -f $(3) .
docker: dockerfiles/Dockerfile.base dockerfiles/Dockerfile.controller dockerfiles/Dockerfile.connector dockercheck
	$(call DOCKER_BUILD,privacyai,latest,dockerfiles/Dockerfile.base)
	$(call DOCKER_BUILD,controller,latest,dockerfiles/Dockerfile.controller)
	$(call DOCKER_BUILD,connector,latest,dockerfiles/Dockerfile.connector)

.PHONY: docker

# ###############################################
# Releasing Docker Images
#
# Using the docker build infrastructure, this section is responsible for
# authenticating to docker hub and pushing built docker containers up with the
# appropriate tags.
# ###############################################
DOCKER_TAG=docker tag dropoutlabs/$(1):$(2) docker.pkg.github.com/dropoutlabs/privacyai/$(1):$(3)
DOCKER_PUSH=docker push docker.pkg.github.com/dropoutlabs/privacyai/$(1):$(2)

docker-tag: dockercheck
	$(call DOCKER_TAG,privacyai,latest,$(VERSION))
	$(call DOCKER_TAG,controller,latest,$(VERSION))
	$(call DOCKER_TAG,connector,latest,$(VERSION))

docker-push-tag: dockercheck
	$(call DOCKER_PUSH,privacyai,$(VERSION))
	$(call DOCKER_PUSH,controller,$(VERSION))
	$(call DOCKER_PUSH,connector,$(VERSION))

docker-push-latest: dockercheck
	$(call DOCKER_PUSH,privacyai,latest)
	$(call DOCKER_PUSH,controller,latest)
	$(call DOCKER_PUSH,connector,latest)

.PHONY: docker-login docker-push-latest docker-push-tag docker-tag

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

releasecheck:
ifneq (yes,$(RELEASE_CONFIRM))
	$(error "Set RELEASE_CONFIRM=yes to really build and push release artifacts")
endif

helm-version-check:
ifeq (,$(shell grep -e $(VERSION) charts/connector/Chart.yaml 2> /dev/null))
	$(error "Version specified in charts/connector/Chart.yaml does not match $(VERSION)")
endif
ifeq (,$(shell grep -e $(VERSION) charts/controller/Chart.yaml 2> /dev/null))
	$(error "Version specified in charts/controller/Chart.yaml does not match $(VERSION)")
endif

publish: releasecheck helmcheck gsutilcheck helm-version-check
	mkdir -p local-dir
	gsutil rsync -d gs://dropout-helm-repo local-dir/
	helm package charts/connector
	helm package charts/controller
	cp *.tgz local-dir
	helm repo index local-dir/ --url https://dropout-helm-repo.storage.googleapis.com
	gsutil rsync -d local-dir/ gs://dropout-helm-repo

.PHONY: publish
