# The `isdubad.toml` configuration file

The `isdubad` server is configured with a [TOML v1.0.0](https://toml.io/en/v1.0.0) file.
Some of the configurables can be overwritten by environment variables. See a full list [here](#env_vars).

An example file can be found [here](./isdubad.toml) (with the default values as comments).

## <a name="env_vars"></a>Environment variables.

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
