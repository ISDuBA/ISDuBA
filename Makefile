# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

.PHONY: all build_isdubad build_importer build_pkg test build_client

all: build_isdubad build_importer build_client test

# See comment here (2024-11-15)
# https://github.com/gocsaf/csaf/blob/3093f717817b9369d390e56d1012eaedcfa19e32/Makefile#L40-L49
GITDESC := $(shell git describe --tags --always)
GITDESCPATCH := $(shell echo '$(GITDESC)' | sed -E 's/v?[0-9]+\.[0-9]+\.([0-9]+)[-+]?.*/\1/')
SEMVERPATCH := $(shell echo $$(( $(GITDESCPATCH) + 1 )))
# Hint: The second regexp in the next line only matches
#       if there is a hyphen (`-`) followed by a number,
#       by which we assume that git describe has added a string after the tag
SEMVER := $(shell echo '$(GITDESC)' | sed -E -e 's/^v//' -e 's/([0-9]+\.[0-9]+\.)([0-9]+)(-[1-9].*)/\1$(SEMVERPATCH)\3/' )
testsemver:
	@echo from \'$(GITDESC)\' transformed to \'$(SEMVER)\'

LDFLAGS=-ldflags "-X github.com/ISDuBA/ISDuBA/pkg/version.SemVersion=$(SEMVER)"
GO_FLAGS=$(LDFLAGS)

# Build for coverage profile generation
ifeq ($(BUILD_COVER), true)
GO_FLAGS += "-cover"
endif


build_importer: build_pkg
	cd cmd/bulkimport && go build $(GO_FLAGS)

build_isdubad: build_pkg
	cd cmd/isdubad && go build $(GO_FLAGS)

build_pkg:
	cd pkg && go build $(GO_FLAGS) ./...

build_client:
	cd client && npm install && npm run build

test:
	go test ./...

DISTNAME := isduba-$(SEMVER)
DISTDIR := dist/$(DISTNAME)
dist: build_isdubad build_client
	mkdir -p $(DISTDIR)
	cp cmd/isdubad/isdubad $(DISTDIR)/
	mkdir -p $(DISTDIR)/web
	cp -r web/* $(DISTDIR)/web
	mkdir -p $(DISTDIR)/docs
	cp -r docs/*.md $(DISTDIR)/docs
	cp -r docs/images/*.svg $(DISTDIR)/docs
	cp -r docs/*.toml $(DISTDIR)/docs
	tar -cvmlzf $(DISTNAME)-gnulinux-amd64.tar.gz $(DISTDIR)
