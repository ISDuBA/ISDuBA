// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package config implements the configuration mechanisms.
package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ISDuBA/ISDuBA/pkg/ginkeycloak"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// DefaultConfigFile is the name of the default config file.
const DefaultConfigFile = "isduba.toml"

const (
	defaultAdvisoryUploadLimit   = 512 * 1024 * 1024
	defaultAnonymousEventLogging = false
)

const (
	defaultLogFile   = "isduba.log"
	defaultLogLevel  = slog.LevelInfo
	defaultLogSource = false
	defaultLogJSON   = false
)

const (
	defaultKeycloakURL           = "http://localhost:8080"
	defaultKeycloakRealm         = "isduba"
	defaultKeycloakCertsCaching  = 8 * time.Hour
	defaultKeycloakTimeout       = 30 * time.Second
	defaultKeycloakFullCertsPath = ""
)

const (
	defaultWebHost    = "localhost"
	defaultWebPort    = 8081
	defaultWebGinMode = "release"
	defaultWebStatic  = "web"
)

const (
	defaultDatabaseHost                    = "localhost"
	defaultDatabasePort                    = 5432
	defaultDatabaseDatabase                = "isduba"
	defaultDatabaseUser                    = "isduba"
	defaultDatabasePassword                = "isduba"
	defaultDatabaseAdminDatabase           = "postgres"
	defaultDatabaseAdminUser               = "postgres"
	defaultDatabaseAdminPassword           = "postgres"
	defaultDatabaseMigrate                 = false
	defaultDatabaseTerminateAfterMigration = true
	defaultMaxQueryDuration                = 30 * time.Second
)

var defaultPublishersTLPs = models.PublishersTLPs{
	"*": []models.TLP{models.TLPWhite},
}

const (
	defaultTempStorageFilesTotal = 10
	defaultTempStorageFilesUser  = 2
	defaultTempStorageDuration   = 30 * time.Minute
)

// HumanSize de-serializes sizes from integer strings
// with suffix "k" (1000), "K" (1024), "m", "M", "g", "G".
// With no suffix given bytes are assumed.
type HumanSize int64

// General are the overarching settings.
type General struct {
	AdvisoryUploadLimit   HumanSize `toml:"advisory_upload_limit"`
	AnonymousEventLogging bool      `toml:"anonymous_event_logging"`
}

// Log are the config options for the logging.
type Log struct {
	File   string     `toml:"file"`
	Level  slog.Level `toml:"level"`
	Source bool       `toml:"source"`
	JSON   bool       `toml:"json"`
}

// Keycloak are the config options for Keycloak.
type Keycloak struct {
	URL           string        `toml:"url"`
	Realm         string        `toml:"realm"`
	CertsCaching  time.Duration `toml:"certs_caching"`
	Timeout       time.Duration `toml:"timeout"`
	FullCertsPath string        `toml:"full_certs_path"`
}

// Web are the config options for the web interface.
type Web struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	GinMode string `toml:"gin_mode"`
	Static  string `toml:"static"`
}

// Database are the config options for the database.
type Database struct {
	Host                    string        `toml:"host"`
	Port                    int           `toml:"port"`
	Database                string        `toml:"database"`
	User                    string        `toml:"user"`
	Password                string        `toml:"password"`
	AdminUser               string        `toml:"admin_user"`
	AdminDatabase           string        `toml:"admin_database"`
	AdminPassword           string        `toml:"admin_password"`
	Migrate                 bool          `toml:"migrate"`
	TerminateAfterMigration bool          `toml:"terminate_after_migration"`
	MaxQueryDuration        time.Duration `toml:"max_query_duration"`
}

// TempStore are the config options for the temporary document storage.
type TempStore struct {
	FilesTotal      int           `toml:"files_total"`
	FilesUser       int           `toml:"files_user"`
	StorageDuration time.Duration `toml:"storage_duration"`
}

// Config are all the configuration options.
type Config struct {
	General        General               `toml:"general"`
	Log            Log                   `toml:"log"`
	Keycloak       Keycloak              `toml:"keycloak"`
	Web            Web                   `toml:"web"`
	Database       Database              `toml:"database"`
	PublishersTLPs models.PublishersTLPs `toml:"publishers_tlps"`
	TempStore      TempStore             `toml:"temp_storage"`
}

// URL creates a connection URL from the configured credentials.
func (db *Database) URL() string {
	url := url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(db.User, db.Password),
		Host:   fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path:   db.Database,
	}
	return url.String()
}

// AdminURL creates a connection URL from the configured credentials.
func (db *Database) AdminURL() string {
	url := url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(db.AdminUser, db.AdminPassword),
		Host:   fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path:   db.AdminDatabase,
	}
	return url.String()
}

// AdminUserURL a connection URL from the configured credentials.
func (db *Database) AdminUserURL() string {
	url := url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(db.AdminUser, db.AdminPassword),
		Host:   fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path:   db.Database,
	}
	return url.String()
}

// Addr returns the combined address the web server should bind to.
func (w *Web) Addr() string {
	return fmt.Sprintf("%s:%d", w.Host, w.Port)
}

// Config returns a Keycloak Config configured by the given settings.
func (kc *Keycloak) Config(mapper ginkeycloak.ClaimMapperFunc) *ginkeycloak.Config {
	return ginkeycloak.NewConfig(
		kc.URL,
		kc.Realm,
	).With(
		ginkeycloak.Cache(kc.CertsCaching),
		ginkeycloak.FullCertsPath(kc.FullCertsPath),
		ginkeycloak.Timeout(kc.Timeout),
		ginkeycloak.CustomClaimsMapper(mapper),
	)
}

// Load loads the configuration from a given file. An empty string
// resorts to the default configuration.
func Load(file string) (*Config, error) {
	cfg := &Config{
		General: General{
			AdvisoryUploadLimit:   defaultAdvisoryUploadLimit,
			AnonymousEventLogging: defaultAnonymousEventLogging,
		},
		Log: Log{
			File:   defaultLogFile,
			Level:  defaultLogLevel,
			Source: defaultLogSource,
			JSON:   defaultLogJSON,
		},
		Keycloak: Keycloak{
			URL:           defaultKeycloakURL,
			Realm:         defaultKeycloakRealm,
			CertsCaching:  defaultKeycloakCertsCaching,
			Timeout:       defaultKeycloakTimeout,
			FullCertsPath: defaultKeycloakFullCertsPath,
		},
		Web: Web{
			Host:    defaultWebHost,
			Port:    defaultWebPort,
			GinMode: defaultWebGinMode,
			Static:  defaultWebStatic,
		},
		Database: Database{
			Host:                    defaultDatabaseHost,
			Port:                    defaultDatabasePort,
			Database:                defaultDatabaseDatabase,
			User:                    defaultDatabaseUser,
			Password:                defaultDatabasePassword,
			AdminDatabase:           defaultDatabaseAdminDatabase,
			AdminUser:               defaultDatabaseAdminUser,
			AdminPassword:           defaultDatabaseAdminPassword,
			Migrate:                 defaultDatabaseMigrate,
			TerminateAfterMigration: defaultDatabaseTerminateAfterMigration,
			MaxQueryDuration:        defaultMaxQueryDuration,
		},
		PublishersTLPs: defaultPublishersTLPs,
		TempStore: TempStore{
			FilesTotal:      defaultTempStorageFilesTotal,
			FilesUser:       defaultTempStorageFilesUser,
			StorageDuration: defaultTempStorageDuration,
		},
	}
	if file != "" {
		md, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return nil, err
		}
		// Don't accept unknown entries in config file.
		if undecoded := md.Undecoded(); len(undecoded) != 0 {
			return nil, fmt.Errorf("config: could not parse %q", undecoded)
		}
	}
	if err := cfg.fillFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) fillFromEnv() error {
	var (
		storeString    = store(noparse)
		storeInt       = store(strconv.Atoi)
		storeBool      = store(strconv.ParseBool)
		storeLevel     = store(storeLevel)
		storeDuration  = store(time.ParseDuration)
		storeHumanSize = store(storeHumanSize)
	)
	return storeFromEnv(
		envStore{"ISDUBA_ADVISORY_UPLOAD_LIMIT", storeHumanSize(&cfg.General.AdvisoryUploadLimit)},
		envStore{"ISDUBA_ANONYMOUS_EVENT_LOGGING", storeBool(&cfg.General.AnonymousEventLogging)},
		envStore{"ISDUBA_LOG_FILE", storeString(&cfg.Log.File)},
		envStore{"ISDUBA_LOG_LEVEL", storeLevel(&cfg.Log.Level)},
		envStore{"ISDUBA_LOG_JSON", storeBool(&cfg.Log.JSON)},
		envStore{"ISDUBA_LOG_SOURCE", storeBool(&cfg.Log.Source)},
		envStore{"ISDUBA_KEYCLOAK_URL", storeString(&cfg.Keycloak.URL)},
		envStore{"ISDUBA_KEYCLOAK_REALM", storeString(&cfg.Keycloak.Realm)},
		envStore{"ISDUBA_KEYCLOAK_TIMEOUT", storeDuration(&cfg.Keycloak.Timeout)},
		envStore{"ISDUBA_KEYCLOAK_CERTS_CACHING", storeDuration(&cfg.Keycloak.CertsCaching)},
		envStore{"ISDUBA_KEYCLOAK_FULL_CERTS_PATH", storeString(&cfg.Keycloak.FullCertsPath)},
		envStore{"ISDUBA_WEB_HOST", storeString(&cfg.Web.Host)},
		envStore{"ISDUBA_WEB_PORT", storeInt(&cfg.Web.Port)},
		envStore{"ISDUBA_WEB_GIN_MODE", storeString(&cfg.Web.GinMode)},
		envStore{"ISDUBA_WEB_STATIC", storeString(&cfg.Web.Static)},
		envStore{"ISDUBA_DB_HOST", storeString(&cfg.Database.Host)},
		envStore{"ISDUBA_DB_PORT", storeInt(&cfg.Database.Port)},
		envStore{"ISDUBA_DB_DATABASE", storeString(&cfg.Database.Database)},
		envStore{"ISDUBA_DB_USER", storeString(&cfg.Database.User)},
		envStore{"ISDUBA_DB_PASSWORD", storeString(&cfg.Database.Password)},
		envStore{"ISDUBA_DB_ADMIN_DATABASE", storeString(&cfg.Database.AdminDatabase)},
		envStore{"ISDUBA_DB_ADMIN_USER", storeString(&cfg.Database.AdminUser)},
		envStore{"ISDUBA_DB_ADMIN_PASSWORD", storeString(&cfg.Database.AdminPassword)},
		envStore{"ISDUBA_DB_MIGRATE", storeBool(&cfg.Database.Migrate)},
		envStore{"ISDUBA_DB_TERMINATE_AFTER_MIGRATION", storeBool(&cfg.Database.TerminateAfterMigration)},
		envStore{"ISDUBA_DB_MAX_QUERY_DURATION", storeDuration(&cfg.Database.MaxQueryDuration)},
		envStore{"ISDUBA_TEMP_STORAGE_FILES_TOTAL", storeInt(&cfg.TempStore.FilesTotal)},
		envStore{"ISDUBA_TEMP_STORAGE_FILES_USER", storeInt(&cfg.TempStore.FilesUser)},
		envStore{"ISDUBA_TEMP_STORAGE_DURATION", storeDuration(&cfg.TempStore.StorageDuration)},
	)
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (hs *HumanSize) UnmarshalText(b []byte) error {
	scale := int64(1)
	if l := len(b); l > 0 {
		switch b[l-1] {
		case 'k':
			scale = 1000
		case 'K':
			scale = 1024
		case 'm':
			scale = 1000 * 1000
		case 'M':
			scale = 1024 * 1024
		case 'g':
			scale = 1000 * 1000 * 1000
		case 'G':
			scale = 1024 * 1024 * 1024
		default:
			goto noUnits
		}
		b = b[:l-1]
	}
noUnits:
	x, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*hs = HumanSize(scale * x)
	return nil
}
