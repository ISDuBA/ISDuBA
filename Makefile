# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

.PHONY: all build_isdubad build_importer build_pkg test build_client

all: build_isdubad build_importer build_client test

# See comment here (2024-04-05)
# https://github.com/csaf-poc/csaf_distribution/blob/d909e9de151d5845fe0c0d5b9db2152f9db25e90/Makefile#L40-L49

GITDESC := $(shell git describe --tags --always)
GITDESCPATCH := $(shell echo '$(GITDESC)' | sed -E 's/v?[0-9]+\.[0-9]+\.([0-9]+)[-+]?.*/\1/')
SEMVERPATCH := $(shell echo $$(( $(GITDESCPATCH) + 1 )))
# Hint: The regexp in the next line only matches if there is a hyphen (`-`)
#       followed by a number, by which we assume that git describe
#       has added a string after the tag
SEMVER := $(shell echo '$(GITDESC)' | sed -E 's/v?([0-9]+\.[0-9]+\.)([0-9]+)(-[1-9].*)/\1$(SEMVERPATCH)\3/' )
testsemver:
	@echo from \'$(GITDESC)\' transformed to \'$(SEMVER)\'

LDFLAGS=-ldflags "-X github.com/ISDuBA/ISDuBA/pkg/version.SemVersion=$(SEMVER)"


build_importer: build_pkg
	cd cmd/bulkimport && go build $(LDFLAGS)

build_isdubad: build_pkg
	cd cmd/isdubad && go build $(LDFLAGS)

build_pkg:
	cd pkg && go build $(LDFLAGS) ./...

build_client:
	cd client && npm install && npm run build

test:
	go test ./...

DISTDIR := isduba-$(SEMVER)
dist: build_isdubad build_client
	mkdir -p dist
	cp cmd/isdubad/isdubad dist/
	mkdir -p dist/web
	cp -r web/* dist/web
	cd dist/ ; tar -cvmlzf $(DISTDIR)-gnulinux-amd64.tar.gz *
