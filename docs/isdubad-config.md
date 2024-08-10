<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# The `isdubad.toml` configuration file

The `isdubad` server is configured with a [TOML v1.0.0](https://toml.io/en/v1.0.0) file.
Some of the configurables can be overwritten by environment variables. See a full list [here](#env_vars).

An example file can be found [here](./example_isdubad.toml) (with the default values as comments).

## Sections

The configuration consists of the following sections:

- [`[general]`](#section_general) General parameters
- [`[log]`](#section_log) Logging
- [`[keycloak]`](#section_keycloak) Keycloak
- [`[web]`](#section_web) Web interface
- [`[database]`](#section_database) Database credentials
- [`[publishers_tlps]`](#section_publishers_tlps) publishers/TLPs filters
- [`[temp_storage]`](#section_temp_storage) temporary document storage
- [`[sources]`](#section_sources) Sources
- [`[client]`](#section_client) Client configuration

### <a name="section_general"></a> Section `[general]` General parameters

- `advisory_upload_limit`: Limits the size of a CSAF document to be uploaded.
  Defaults to `"512K"`. Recognized unit suffixes are
  `k`/`K` for 1000/1024, `m`/`M` for 1000<sup>2</sup>/1024<sup>2</sup>,
  `g`/`G` 1000<sup>3</sup>/1024<sup>3</sup> and none for bytes.
- `anonymous_event_logging`: Indicates that the event logging of the document
  workflow life cycle should be stored with no user. Defaults to `false`.

### <a name="section_log"></a> Section `[log]` Logging

- `file`: File to log to. An empty string logs to stderr. Defaults to `"isduba.log"`.
- `level`: Log level. Possible values are `"debug"`, `"info"`, `"warn"` and `"error"`. Defaults to `"info"`.
- `source`: Add source reference to log output. Defaults to `false`.
- `json`: Log as JSON lines. Defaults to `false`.

### <a name="section_keycloak"></a> Section `[keycloak]` Keycloak

- `url`: Defaults to `"http://localhost:8080"`.
- `realm`: Name of the realm used be the server. Defaults to `"isduba"`.
- `certs_caching`: How long should signing certificats from the Keycloak should be cached before reasked. Defaults to `"8h"`.
- `timeout`: How long should we wait for reactions from the Keycloak server. Defaults to `"30s"`.
- `full_certs_path`: Special URL to fetch the signing certificats from. Defaults to `""`.

### <a name="section_web"></a> Section `[web]` Web interface

- `host`: Interface the web server listens on. Defaults to `"localhost"`.
- `port`: Port the web server listens on. Defaults to `8081`.
- `gin_mode`: Mode the Gin middleware is running in. Defaults to `"release"`.
- `static`: Folder to be served under **<http://host:port/>**. Defaults to `"web"`.

### <a name="section_database"></a> Section `[database]` Database credentials

- `host`: Host of the database server. Defaults to `"localhost"`.
- `port`: Port of the database server. Defaults to `5432`.
- `database`: Name of the database. Defaults to `"isduba"`.
- `user`: Name of the database user. Defaults to `"isduba"`.
- `password`: Passwordof the database user. Defaults to `"isduba"`.
- `admin_user`: Name of an admin database user. Only needed when migration is needed. Defaults to `"postgres"`.
- `admin_database`: Name of an admin database. Only needed in case of migrations. Defaults to `"postgres"`.
- `admin_password`: Password of the admin user. For migrations only. Defaults to `"postgres"`.
- `migrate`: Should a migration be performed if needed? Better triggered by the **ISDUBA_DB_MIGRATE** env variable. Defaults to `false`.
- `terminate_after_migration` When a migration is started the program terminates by default.
- `max_query_duration`: How long a user provided database may last at max. Defaults to `"30s"`.

### <a name="section_publishers_tlps"></a> Section `[publishers_tlps]` publishers/TLP filters

Is a table of pairs of `publisher` and `tlps` of default access.
Defaults to:

```
[publishers_tlps]
'*' = ["WHITE"]
```

An empty `publisher` means all not explicity stated. `publisher`s with non empty values have a higher priority.
Valid values for `tlps` are the [Traffic Light Protocol](https://en.wikipedia.org/wiki/Traffic_Light_Protocol) 1 values
`WHITE`, `GREEN`, `AMBER` and `RED`.

### <a name="section_temp_storage"></a> Section `[temp_storage]` Temporary document storage

- `files_total`: Max number of files hold in temp storage. Defaults to `10`.
- `files_user`: Max number of files hold in temp storage per user. Defaults to `2`.
- `storage_duration`: Ensured storage duration in temp storage. Defaults to `"30m"` minutes.

### <a name="section_sources"></a> Section `[sources]` Sources

- `download_slots`: The number of concurrent downloads from the sources. Defaults to `100`.
- `max_slots_per_source`: The number of concurrent downloads per source. Defaults to `2`.
- `max_rate_per_source`: The Number of requests per source per second. Defaults to `0` (unlimited).
- `feed_refresh`: Duration between re-asking source for a new updated feed index. Defaults to `15m`.
- `feed_log_level`: The log level per feed. Valid values are `debug`, `info`, `warn`, `error`. Defaults to `info`.
- `feed_importer`: Name of the user that is doing the feed imports. Defaults to `feedimporter`.
- `publishers_tlps`: Rules what the feed import is allowed to import. Defaults to `{ "*" = [ "WHITE", "GREEN", "AMBER", "RED" ] }`
- `default_message`: The message that should be displayed inside the source manager.

### <a name="section_client"></a> Section `[client]` Client configuration

- `keycloak_url`: The URL where the Keycloak server is located. Defaults to same as `keycloak.url`.
- `keycloak_realm`: The name of the Keycloak realm. Defaults to "isduba".
- `keycloak_client_id`: The public client identifier. Defaults to "auth".
- `update_interval`: Specifies how often the token should be renewed. Defaults to "5m".
- `application_uri`: The base URL of the application. Defaults to `http://localhost:8081` where **localhost** is the same as `web.host` and **8081** is the same as `web.port`.
- `idle_timeout`: When the user should be logged out after inactivity. Defaults to "30m".

## <a name="env_vars"></a>Environment variables

| Env variable                          | Overwrites                           |
| ------------------------------------- | ------------------------------------ |
| `ISDUBA_ADVISORY_UPLOAD_LIMIT`        | `general advisory_upload_limit`      |
| `ISDUBA_ANONYMOUS_EVENT_LOGGING`      | `general anonymous_event_logging`    |
| `ISDUBA_LOG_FILE`                     | `log file`                           |
| `ISDUBA_LOG_LEVEL`                    | `log level`                          |
| `ISDUBA_LOG_JSON"`                    | `log json`                           |
| `ISDUBA_LOG_SOURCE`                   | `log source`                         |
| `ISDUBA_KEYCLOAK_URL`                 | `keycloak url`                       |
| `ISDUBA_KEYCLOAK_REALM`               | `keycloak realm`                     |
| `ISDUBA_KEYCLOAK_TIMEOUT`             | `keycloak timeout`                   |
| `ISDUBA_KEYCLOAK_CERTS_CACHING`       | `keycloak certs_caching`             |
| `ISDUBA_KEYCLOAK_FULL_CERTS_PATH`     | `keycloak full_certs_path`           |
| `ISDUBA_WEB_HOST`                     | `web host`                           |
| `ISDUBA_WEB_PORT`                     | `web port`                           |
| `ISDUBA_WEB_GIN_MODE`                 | `web gin_mode`                       |
| `ISDUBA_WEB_STATIC`                   | `web static`                         |
| `ISDUBA_DB_HOST`                      | `database host`                      |
| `ISDUBA_DB_PORT`                      | `database port`                      |
| `ISDUBA_DB_DATABASE`                  | `database database`                  |
| `ISDUBA_DB_USER`                      | `database user`                      |
| `ISDUBA_DB_PASSWORD`                  | `database password`                  |
| `ISDUBA_DB_ADMIN_DATABASE`            | `database admin_database`            |
| `ISDUBA_DB_ADMIN_USER`                | `database admin_user`                |
| `ISDUBA_DB_ADMIN_PASSWORD`            | `database admin_password`            |
| `ISDUBA_DB_MIGRATE`                   | `database migrate`                   |
| `ISDUBA_DB_TERMINATE_AFTER_MIGRATION` | `database terminate_after_migration` |
| `ISDUBA_DB_MAX_QUERY_DURATION`        | `database max_query_duration`        |
| `ISDUBA_TEMP_STORAGE_FILES_TOTAL`     | `temp_storage files_total`           |
| `ISDUBA_TEMP_STORAGE_FILES_USER`      | `temp_storage files_user`            |
| `ISDUBA_TEMP_STORAGE_DURATION`        | `temp_storage storage_duration`      |
| `ISDUBA_SOURCES_DOWNLOAD_SLOTS`       | `sources download_slots`             |
| `ISDUBA_SOURCES_MAX_SLOTS_PER_SOURCE` | `sources max_slots_per_source`       |
| `ISDUBA_SOURCES_MAX_RATE_PER_SOURCE`  | `sources max_rate_per_source`        |
| `ISDUBA_SOURCES_FEED_REFRESH`         | `sources feed_refresh`               |
| `ISDUBA_SOURCES_FEED_LOG_LEVEL`       | `sources feed_log_level`             |
| `ISDUBA_SOURCES_FEED_IMPORTER`        | `sources feed_importer`              |
| `ISDUBA_SOURCES_DEFAULT_MESSAGE`      | `sources default_message`            |
| `ISDUBA_CLIENT_KEYCLOAK_URL`          | `client keycloak_url`                |
| `ISDUBA_CLIENT_KEYCLOAK_REALM`        | `client keycloak_realm`              |
| `ISDUBA_CLIENT_KEYCLOAK_CLIENT_ID`    | `client keycloak_client_id`          |
| `ISDUBA_CLIENT_UPDATE_INTERVAL`       | `client update_interval`             |
| `ISDUBA_CLIENT_APPLICATION_URI`       | `client application_uri`             |
| `ISDUBA_CLIENT_IDLE_TIMEOUT`          | `client idle_timeout`                |
