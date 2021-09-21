#############################
# Global vars
#############################
PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  ?= $(shell git describe --tags --always --dirty | sed -e '/^v/s/^v\(.*\)$$/\1/g')
# Last released version (not dirty)
PROJECT_VER_TAGGED  := $(shell git describe --tags --always --abbrev=0 | sed -e '/^v/s/^v\(.*\)$$/\1/g')

SRCDIR       ?= .
GO            = go

# The root module (from go.mod)
PROJECT_MODULE  ?= $(shell $(GO) list -m)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

#############################
# Targets
#############################
all: build

# Humans running make:
build: git-hooks check-version clean lint test cover-report compile

# Build command for CI tooling
build-ci: check-version clean lint test compile-only

# All clean commands
clean: cover-clean compile-clean release-clean

# Import fragments
include build/compile.mk
include build/deps.mk
include build/docker.mk
include build/document.mk
include build/generate.mk
include build/lint.mk
include build/release.mk
include build/test.mk
include build/tools.mk
include build/util.mk

.PHONY: all build build-ci clean
