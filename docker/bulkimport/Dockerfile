# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

FROM library/golang:1.22.1-alpine as build
RUN apk update && apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile

ARG BUILD_VERSION
RUN make GITDESC=$BUILD_VERSION build_importer


FROM library/golang:1.22.1-alpine as importer

COPY --from=build /app/cmd/bulkimport/bulkimport /app/bulkimport

WORKDIR /app

CMD ["./bulkimport", "-database", "isduba", "-user", "isduba", "-password", "isduba", "-host", "isduba-db", "/data"]

