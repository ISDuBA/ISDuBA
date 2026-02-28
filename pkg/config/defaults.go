// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package config

import (
	"log/slog"
	"time"

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
	defaultSourcesChecking       = 2 * time.Hour
	defaultKeepFeedLogs          = 3 * 31 * 24 * time.Hour
)

const (
	defaultForwarderUpdateInterval = 5 * time.Minute
	defaultForwarderStratgy        = ForwarderStrategyAll
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
