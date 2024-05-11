# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

FROM library/node:20.10-alpine as client
RUN apk update && apk add --no-cache make

WORKDIR /app

ARG PUBLIC_KEYCLOAK_URL
ARG PUBLIC_KEYCLOAK_REALM
ARG PUBLIC_KEYCLOAK_CLIENTID

ENV PUBLIC_KEYCLOAK_URL=$PUBLIC_KEYCLOAK_URL
ENV PUBLIC_KEYCLOAK_REALM=$PUBLIC_KEYCLOAK_REALM
ENV PUBLIC_KEYCLOAK_CLIENTID=$PUBLIC_KEYCLOAK_CLIENTID

COPY client/package.json client/package.json
RUN cd client && npm install

COPY client client
COPY Makefile Makefile
RUN make build_client


FROM library/golang:1.22.1-alpine as server
RUN apk update && apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile

ARG BUILD_VERSION
RUN make GITDESC=$BUILD_VERSION build_isdubad
ENV ISDUBA_DB_MIGRATE true


FROM scratch

COPY --from=server /app/cmd/isdubad/isdubad /app/isdubad
COPY --from=client /app/client/build /app/web
COPY ./docker/server/isdubad.toml /app/isdubad.toml

WORKDIR /app

CMD ["./isdubad", "-c", "./isdubad.toml"]

