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
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gocsaf/csaf/v3/csaf"

	"github.com/ISDuBA/ISDuBA/pkg/ginkeycloak"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// DefaultConfigFile is the name of the default config file.
const DefaultConfigFile = "isduba.toml"

const (
	defaultAdvisoryUploadLimit   = 512 * 1024 * 1024
	defaultAnonymousEventLogging = false
)

var (
	defaultURLPorts      = []PortRange{{80, 80}, {443, 443}}
	defaultBlockedRanges = []string{
		// Taken from https://gist.github.com/stefansundin/32e8399f0c67c07c372b5ab51560e004
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	}
)

const (
	defaultBlockLoopback = true
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

var (
	defaultPublishersTLPs = models.PublishersTLPs{
		"*": []models.TLP{models.TLPWhite},
	}
	defaultSourcesPublishersTLPs = models.PublishersTLPs{
		"*": []models.TLP{
			models.TLPWhite,
			models.TLPGreen,
			models.TLPAmber,
			models.TLPRed,
		},
	}
)

const (
	defaultTempStorageFilesTotal = 10
	defaultTempStorageFilesUser  = 2
	defaultTempStorageDuration   = 30 * time.Minute
)

const (
	defaultSourcesDownloadSlots     = 100
	defaultSourcesMaxSlotsPerSource = 2
	defaultSourcesMaxRatePerSlot    = 0
	defaultSourcesOpenPGPCaching    = 24 * time.Hour
	defaultSourcesFeedRefresh       = 15 * time.Minute
	defaultSourcesTimeout           = 30 * time.Second
	defaultSourcesFeedLogLevel      = InfoFeedLogLevel
	defaultSourcesFeedImporter      = "feedimporter"
	defaultSourcesDefaultMessage    = "Missing something? To suggest new CSAF sources, " +
		"please contact your CSAF source manager. Otherwise contact your administrator."
	defaultSourcesStrictMode     = true
	defaultSourcesSecure         = true
	defaultSourcesSignatureCheck = true
	defaultSourcesAge            = 17520 * time.Hour
	defaultSourcesAESKey         = ""
	defaultSourcesChecking       = 2 * time.Hour
	defaultKeepFeedLogs          = 3 * 31 * 24 * time.Hour
)

const (
	defaultForwarderUpdateInterval = 5 * time.Minute
)

const (
	defaultRemoteValidatorURL   = ""
	defaultRemoteValidatorCache = ""
)

var defaultRemoteValidatorPresets = []string{"mandatory"}

const (
	defaultClientKeycloakRealm    = "isduba"
	defaultClientKeycloakClientID = "auth"
	defaultClientUpdateInterval   = 5 * time.Minute
	defaultClientIdleTimeout      = 30 * time.Minute
)

const (
	defaultAggregatorsTimeout        = 30 * time.Second
	defaultAggregatorsUpdateInterval = 1 * time.Hour
)

// HumanSize de-serializes sizes from integer strings
// with suffix "k" (1000), "K" (1024), "m", "M", "g", "G".
// With no suffix given bytes are assumed.
type HumanSize int64

// General are the overarching settings.
type General struct {
	AdvisoryUploadLimit   HumanSize   `toml:"advisory_upload_limit"`
	AnonymousEventLogging bool        `toml:"anonymous_event_logging"`
	AllowedPorts          []PortRange `toml:"allowed_ports"`
	BlockLoopback         bool        `toml:"block_loopback"`
	BlockedRanges         []IPRange   `toml:"blocked_ranges"`
	AllowedIPs            []net.IP    `toml:"allowed_ips"`
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

// Sources are the config options for downloading sources.
type Sources struct {
	DownloadSlots     int                   `toml:"download_slots"`
	MaxSlotsPerSource int                   `toml:"max_slots_per_source"`
	MaxRatePerSource  float64               `toml:"max_rate_per_source"`
	OpenPGPCaching    time.Duration         `toml:"openpgp_caching"`
	FeedRefresh       time.Duration         `toml:"feed_refresh"`
	Timeout           time.Duration         `toml:"timeout"`
	FeedLogLevel      FeedLogLevel          `tomt:"feed_log_level"`
	PublishersTLPs    models.PublishersTLPs `toml:"publishers_tlps"`
	FeedImporter      string                `toml:"feed_importer"`
	DefaultMessage    string                `toml:"default_message"`
	StrictMode        bool                  `toml:"strict_mode"`
	Secure            bool                  `toml:"secure"`
	SignatureCheck    bool                  `toml:"signature_check"`
	DefaultAge        time.Duration         `toml:"default_age"`
	AESKey            string                `toml:"aes_key"`
	Checking          time.Duration         `toml:"checking"`
	KeepFeedLogs      time.Duration         `toml:"keep_feed_logs"`
}

// ForwardTarget are the config options for the forward target.
type ForwardTarget struct {
	URL               string        `toml:"url"`
	Name              string        `toml:"name"`
	ClientPrivateCert string        `toml:"private_cert"`
	ClientPublicCert  string        `toml:"public_cert"`
	Publisher         *string       `toml:"publisher"`
	Header            []string      `toml:"header"`
	Automatic         bool          `toml:"automatic"`
	Timeout           time.Duration `toml:"timeout"`
}

// Forwarder are the config options for the document forwarder.
type Forwarder struct {
	Targets        []ForwardTarget `toml:"target"`
	UpdateInterval time.Duration   `toml:"update_interval"`
}

// Aggregators are the config options for the aggregators.
type Aggregators struct {
	Timeout        time.Duration `toml:"timeout"`
	UpdateInterval time.Duration `toml:"update_interval"`
}

// Client are the config options for the client.
type Client struct {
	KeycloakURL      string        `toml:"keycloak_url" json:"keycloak_url"`
	KeycloakRealm    string        `toml:"keycloak_realm" json:"keycloak_realm"`
	KeycloakClientID string        `toml:"keycloak_client_id" json:"keycloak_client_id"`
	UpdateInterval   time.Duration `toml:"update_interval" json:"update_interval" swaggertype:"integer"`
	IdleTimeout      time.Duration `toml:"idle_timeout" json:"idle_timeout" swaggertype:"integer"`
}

// Config are all the configuration options.
type Config struct {
	General         General                     `toml:"general"`
	Log             Log                         `toml:"log"`
	Keycloak        Keycloak                    `toml:"keycloak"`
	Web             Web                         `toml:"web"`
	Database        Database                    `toml:"database"`
	PublishersTLPs  models.PublishersTLPs       `toml:"publishers_tlps"`
	TempStore       TempStore                   `toml:"temp_storage"`
	Sources         Sources                     `toml:"sources"`
	RemoteValidator csaf.RemoteValidatorOptions `toml:"remote_validator"`
	Client          Client                      `toml:"client"`
	Forwarder       Forwarder                   `toml:"forwarder"`
	Aggregators     Aggregators                 `toml:"aggregators"`
}

func escape(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `\'`) + `'`
}

func connString(user, password, host string, port int, database string) string {
	return fmt.Sprintf("user=%s "+
		"password=%s "+
		"host=%s "+
		"port=%d "+
		"database=%s",
		escape(user),
		escape(password),
		escape(host),
		port,
		escape(database))
}

// ConnString creates a connection ConnString from the configured credentials.
func (db *Database) ConnString() string {
	return connString(db.User, db.Password, db.Host, db.Port, db.Database)
}

// AdminConnString creates a connection URL from the configured credentials.
func (db *Database) AdminConnString() string {
	return connString(db.AdminUser, db.AdminPassword, db.Host, db.Port, db.AdminDatabase)
}

// AdminUserConnString a connection URL from the configured credentials.
func (db *Database) AdminUserConnString() string {
	return connString(db.AdminUser, db.AdminPassword, db.Host, db.Port, db.Database)
}

// Addr returns the combined address the web server should bind to.
func (w *Web) Addr() string {
	return net.JoinHostPort(w.Host, strconv.Itoa(w.Port))
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
			AllowedPorts:          nil,
			BlockLoopback:         defaultBlockLoopback,
			BlockedRanges:         nil,
			AllowedIPs:            nil,
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
		Sources: Sources{
			DownloadSlots:     defaultSourcesDownloadSlots,
			MaxSlotsPerSource: defaultSourcesMaxSlotsPerSource,
			MaxRatePerSource:  defaultSourcesMaxRatePerSlot,
			OpenPGPCaching:    defaultSourcesOpenPGPCaching,
			FeedRefresh:       defaultSourcesFeedRefresh,
			Timeout:           defaultSourcesTimeout,
			FeedLogLevel:      defaultSourcesFeedLogLevel,
			FeedImporter:      defaultSourcesFeedImporter,
			PublishersTLPs:    defaultSourcesPublishersTLPs,
			DefaultMessage:    defaultSourcesDefaultMessage,
			StrictMode:        defaultSourcesStrictMode,
			Secure:            defaultSourcesSecure,
			SignatureCheck:    defaultSourcesSignatureCheck,
			DefaultAge:        defaultSourcesAge,
			Checking:          defaultSourcesChecking,
			KeepFeedLogs:      defaultKeepFeedLogs,
		},
		Forwarder: Forwarder{
			UpdateInterval: defaultForwarderUpdateInterval,
		},
		RemoteValidator: csaf.RemoteValidatorOptions{
			URL:     defaultRemoteValidatorURL,
			Presets: defaultRemoteValidatorPresets,
			Cache:   defaultRemoteValidatorCache,
		},
		Client: Client{
			KeycloakRealm:    defaultClientKeycloakRealm,
			KeycloakClientID: defaultClientKeycloakClientID,
			UpdateInterval:   defaultClientUpdateInterval,
			IdleTimeout:      defaultClientIdleTimeout,
		},
		Aggregators: Aggregators{
			Timeout:        defaultAggregatorsTimeout,
			UpdateInterval: defaultAggregatorsUpdateInterval,
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
	cfg.presetEmptyDefaults()
	return cfg, nil
}

func parsedDefaultBlockedRanges() []IPRange {
	brs := make([]IPRange, 0, len(defaultBlockedRanges))
	for _, cidr := range defaultBlockedRanges {
		_, blocked, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Sprintf("invalid CIDR %q", cidr))
		}
		brs = append(brs, IPRange{IPNet: blocked})
	}
	return brs
}

func (cfg *Config) presetEmptyDefaults() {
	if cfg.General.AllowedPorts == nil {
		cfg.General.AllowedPorts = defaultURLPorts
	}
	if cfg.General.BlockedRanges == nil {
		cfg.General.BlockedRanges = parsedDefaultBlockedRanges()
	}
	if cfg.Client.KeycloakURL == "" {
		cfg.Client.KeycloakURL = cfg.Keycloak.URL
	}
}

func (cfg *Config) fillFromEnv() error {
	var (
		storeString       = store(noparse)
		storeInt          = store(strconv.Atoi)
		storeBool         = store(strconv.ParseBool)
		storeLevel        = store(storeLevel)
		storeDuration     = store(time.ParseDuration)
		storeHumanSize    = store(storeHumanSize)
		storeFeedLogLevel = store(storeFeedLogLevel)
		storeFloat64      = store(parseFloat64)
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
		envStore{"ISDUBA_SOURCES_DOWNLOAD_SLOTS", storeInt(&cfg.Sources.DownloadSlots)},
		envStore{"ISDUBA_SOURCES_MAX_SLOTS_PER_SOURCE", storeInt(&cfg.Sources.MaxSlotsPerSource)},
		envStore{"ISDUBA_SOURCES_MAX_RATE_PER_SOURCE", storeFloat64(&cfg.Sources.MaxRatePerSource)},
		envStore{"ISDUBA_SOURCES_OPENPGP_CACHING", storeDuration(&cfg.Sources.OpenPGPCaching)},
		envStore{"ISDUBA_SOURCES_FEED_REFRESH", storeDuration(&cfg.Sources.FeedRefresh)},
		envStore{"ISDUBA_SOURCES_FEED_LOG_LEVEL", storeFeedLogLevel(&cfg.Sources.FeedLogLevel)},
		envStore{"ISDUBA_SOURCES_FEED_IMPORTER", storeString(&cfg.Sources.FeedImporter)},
		envStore{"ISDUBA_SOURCES_DEFAULT_MESSAGE", storeString(&cfg.Sources.DefaultMessage)},
		envStore{"ISDUBA_SOURCES_STRICT_MODE", storeBool(&cfg.Sources.StrictMode)},
		envStore{"ISDUBA_SOURCES_SECURE", storeBool(&cfg.Sources.Secure)},
		envStore{"ISDUBA_SOURCES_SIGNATURE_CHECK", storeBool(&cfg.Sources.SignatureCheck)},
		envStore{"ISDUBA_SOURCES_TIMEOUT", storeDuration(&cfg.Sources.Timeout)},
		envStore{"ISDUBA_SOURCES_DEFAULT_AGE", storeDuration(&cfg.Sources.DefaultAge)},
		envStore{"ISDUBA_SOURCES_AES_KEY", storeString(&cfg.Sources.AESKey)},
		envStore{"ISDUBA_SOURCES_CHECKING", storeDuration(&cfg.Sources.Checking)},
		envStore{"ISDUBA_SOURCES_KEEP_FEED_LOGS", storeDuration(&cfg.Sources.KeepFeedLogs)},
		envStore{"ISDUBA_REMOTE_VALIDATOR_URL", storeString(&cfg.RemoteValidator.URL)},
		envStore{"ISDUBA_REMOTE_VALIDATOR_CACHE", storeString(&cfg.RemoteValidator.Cache)},
		envStore{"ISDUBA_CLIENT_KEYCLOAK_URL", storeString(&cfg.Client.KeycloakURL)},
		envStore{"ISDUBA_CLIENT_KEYCLOAK_REALM", storeString(&cfg.Client.KeycloakRealm)},
		envStore{"ISDUBA_CLIENT_KEYCLOAK_CLIENT_ID", storeString(&cfg.Client.KeycloakClientID)},
		envStore{"ISDUBA_CLIENT_UPDATE_INTERVAL", storeDuration(&cfg.Client.UpdateInterval)},
		envStore{"ISDUBA_CLIENT_IDLE_TIMEOUT", storeDuration(&cfg.Client.IdleTimeout)},
		envStore{"ISDUBA_FORWARDER_UPDATE_INTERVAL", storeDuration(&cfg.Forwarder.UpdateInterval)},
		envStore{"ISDUBA_AGGREGATORS_TIMEOUT", storeDuration(&cfg.Aggregators.Timeout)},
		envStore{"ISDUBA_AGGREGATORS_UPDATE_INTERVAL", storeDuration(&cfg.Aggregators.UpdateInterval)},
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
