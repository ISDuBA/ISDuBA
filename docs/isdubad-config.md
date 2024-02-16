# The `isdubad.toml` configuration file

The `isdubad` server is configured with a [TOML v1.0.0](https://toml.io/en/v1.0.0) file.
Some of the configurables can be overwritten by environment variables. See a full list [here](#env_vars).

An example file can be found [here](./isdubad.toml) (with the default values as comments).

## Sections

The configuration consists of the following sections:

- [`[general]`](#section_general) General parameters
- [`[log]`](#section_log) Logging
- [`[keycloak]`](#section_keycloak) Keycloak
- [`[web]`](#section_web) Web interface
- [`[database]`](#section_database) Database credentials
- [`[[publishers_tlps]]`](#section_publishers_tlps) publishers/TLPs filters

### <a name="section_general"></a> Section `[general]` General password

- `advisory_upload_limit`: Limits the size of a CSAF document to be uploaded.
   Defaults to `"512K"`. Recognized unit suffixes are
   k`/`K` for 1000/1024, `m`/`M` for 1000<sup>2</sup>/1024<sup>2</sup>,
   `g`/`G` 1000<sup>3</sup>/1024<sup>3</sup> and none for bytes.
- `anonymous_event_logging': Indicates that the event logging of the document
   workflow life cycle should be stored with no user. Defaults to `false`.

### <a name="section_log"></a> Section `[log]` Logging

**TBD**

### <a name="section_keycloak"></a> Section `[keycloak]` Keycloak

**TBD**

### <a name="section_web"></a> Section `[web]` Web interface

**TBD**

### <a name="section_database"></a> Section `[database]` Database credentials

**TBD**

### <a name="section_publishers_tlps"></a> Section `[[publishers_tlps]]` publishers/TLP filters

**TBD**

## <a name="env_vars"></a>Environment variables

| Env variable | Overwrites |
| ------------ | ---------- |
| `ISDUBA_ADVISORY_UPLOAD_LIMIT` | `general advisory_upload_limit` |
| `ISDUBA_ANONYMOUS_EVENT_LOGGING` | `general anonymous_event_logging` |
| `ISDUBA_LOG_FILE` | `log file` |
| `ISDUBA_LOG_LEVEL` | `log level` |
| `ISDUBA_LOG_JSON"` | `log json` |
| `ISDUBA_LOG_SOURCE` | `log source` |
| `ISDUBA_KEYCLOAK_URL` | `keycloak url` |
| `ISDUBA_KEYCLOAK_REALM` | `keycloak realm` |
| `ISDUBA_KEYCLOAK_TIMEOUT` | `keycloak timeout` |
| `ISDUBA_KEYCLOAK_CERTS_CACHING` | `keycloak certs_caching` |
| `ISDUBA_KEYCLOAK_FULL_CERTS_PATH` | `keycloak full_certs_path` |
| `ISDUBA_WEB_HOST` | `web host` |
| `ISDUBA_WEB_PORT` | `web port` |
| `ISDUBA_WEB_GIN_MODE` | `web gin_mode` |
| `ISDUBA_WEB_STATIC` | `web static` |
| `ISDUBA_DB_HOST` | `database host` |
| `ISDUBA_DB_PORT` | `database port` |
| `ISDUBA_DB_DATABASE` | `database database` |
| `ISDUBA_DB_USER` | `database user` |
| `ISDUBA_DB_PASSWORD` | `database password` |
| `ISDUBA_DB_ADMIN_DATABASE` | `database admin_database ` |
| `ISDUBA_DB_ADMIN_USER` | `database admin_user` |
| `ISDUBA_DB_ADMIN_PASSWORD` | `database admin_password` |
| `ISDUBA_DB_MIGRATE` | `database migrate` |
