			######################################################
#
# Variable definitions
#
######################################################

basedir := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
parentdir := $(dir $(patsubst %/,%,$(dir $(basedir))))
gofiles := $(shell find $(basedir) -type f -name '*.go' -not -path "./vendor/*")
dockerenv := $(filter-out export, $(shell echo $(shell minikube docker-env 2> /dev/null)))
timestamp := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

######################################################
#
# Variables which can be overridden
#
######################################################

CODEVERSION := latest
PERSONALITY := $(notdir $(patsubst %/,%,$(dir $(basedir))))
GOPATH := $(parentdir)gopath
PROSIMOBINDIR := $(GOPATH)/bin/prosimo
BINFILE := $(PROSIMOBINDIR)/$(PERSONALITY)
BINDIR := $(BINFILE)-all
PLGDIR := $(BINFILE)-plg
GO := GOPATH=$(GOPATH) GOPRIVATE=git.prosimo.io/prosimoio  GOOS=darwin GOARCH=amd64 go
GONATIVE := GOPATH=$(GOPATH) GOPRIVATE=git.prosimo.io/prosimoio  go
GOIMPORTS := $(GOPATH)/bin/goimports
GOLINT := $(GOPATH)/bin/golint
GITHASH := $(shell git rev-parse --short HEAD)
GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD)

ifeq ($(BUILDENV),production)
	PRODFLAG=TRUE
else
ifeq ($(BUILDENV),staging)
	PRODFLAG=TRUE
else
	PRODFLAG=FALSE
endif
endif

######################################################
#
# Development targets
#
######################################################

.PHONY: version
.PHONY: setup test check check-fmt check-lint check-vet
.PHONY: build image login push all ci-all-gcs ci-build
.PHONY: clean-build clean-build-all clean-push clean-all
.PHONY: build-all install test testacc
.PHONY: stat clean-stat

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=prosimo.com
NAMESPACE=prosimo
NAME=prosimo
BINARY=terraform-provider-${NAME}
MAJOR_VERSION=3
MINOR_VERSION=4
REVISION=4
# BUILDNUMBER=0
VERSION=${MAJOR_VERSION}.$(MINOR_VERSION).$(REVISION)
OS_ARCH=darwin_amd64

default: install

stat:
	cd $(basedir) && tfsec --force-all-dirs

build:
	go fmt prosimo/*
	cd $(basedir) && $(GO) build -o $(BINFILE)

major-version:
	$(shell ./versioning.sh major-version)

minor-version:
	$(shell ./versioning.sh minor-version)

revision:
	$(shell ./versioning.sh revision)	 

.PHONY: docs

docs: tools
		tfplugindocs generate

tools:
	GO111MODULE=on go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	
all: install

install: build
	mkdir -p $(HOME)/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH); \
	cp $(BINFILE) $(HOME)/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

######################################################
#
# Multi OS and Architectures
#
######################################################

GOA := GOPATH=$(GOPATH) GOPRIVATE=git.prosimo.io/prosimoio go

build-all:
	mkdir -p $(BINDIR)/$(VERSION)
	cd $(basedir) && GOOS=darwin GOARCH=amd64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_darwin_amd64
	cd $(basedir) && GOOS=darwin GOARCH=arm64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_darwin_arm64
	cd $(basedir) && GOOS=freebsd GOARCH=386 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_freebsd_386
	cd $(basedir) && GOOS=freebsd GOARCH=amd64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_freebsd_amd64
	cd $(basedir) && GOOS=freebsd GOARCH=arm $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_freebsd_arm
	cd $(basedir) && GOOS=linux GOARCH=386 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_linux_386
	cd $(basedir) && GOOS=linux GOARCH=amd64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_linux_amd64
	cd $(basedir) && GOOS=linux GOARCH=arm $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_linux_arm
	cd $(basedir) && GOOS=openbsd GOARCH=386 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_openbsd_386
	cd $(basedir) && GOOS=openbsd GOARCH=amd64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_openbsd_amd64
	cd $(basedir) && GOOS=solaris GOARCH=amd64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_solaris_amd64
	cd $(basedir) && GOOS=windows GOARCH=386 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_windows_386
	cd $(basedir) && GOOS=windows GOARCH=amd64 $(GOA) build -o $(BINDIR)/$(VERSION)/$(BINARY)_$(VERSION)_windows_amd64

image:
	cd $(basedir) && tar -C $(BINDIR) -cpf $(PERSONALITY).tar .

push:
	gsutil -o GSUtil:parallel_composite_upload_threshold=150M cp $(basedir)$(PERSONALITY).tar gs://$(GCSBUCKET)

clean-build:
	rm -f $(BINFILE)

clean-build-all:
	rm -rf $(BINDIR) && rm -f $(basedir)$(PERSONALITY).tar

clean-push:
	gsutil rm gs://$(GCSBUCKET)$(PERSONALITY).tar

clean-stat:
	rm -f $(PERSONALITY)-static.log

clean-all: clean-push clean-build-all clean-build clean-stat

######################################################
#
# CI targets
#
######################################################
ci-build: build-all
ci-all-gcs: build-all image push stat
