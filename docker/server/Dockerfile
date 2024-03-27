# This file is Free Software under the Apache-2.0 License
# without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
#
# SPDX-License-Identifier: Apache-2.0
#
# SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
# Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

FROM mirror.mobydock.de/library/node:20.10-alpine as client

WORKDIR /app/client

COPY client/package.json .
RUN npm install

COPY client .
COPY ./pkg/version/version.go ../pkg/version/version.go
RUN npm run build


FROM mirror.mobydock.de/library/golang:1.22.1-alpine as server

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd cmd
COPY pkg pkg

WORKDIR /app/cmd/isdubad
RUN go build -v -o /app/isdubad
ENV ISDUBA_DB_MIGRATE true


FROM scratch

COPY --from=server /app/isdubad /app/isdubad
COPY --from=client /app/client/build /app/web
COPY ./docker/server/isdubad.toml /app/isdubad.toml

WORKDIR /app

CMD ["./isdubad", "-c", "./isdubad.toml"]

