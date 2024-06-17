// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package ginkeycloak implements a Gin middleware to handle JWT tokens produced by Keycloak.
package ginkeycloak

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// ClaimMapperFunc is a custom function to map from a JWT token to a Keycloak one.
type ClaimMapperFunc func(func(any) error, *KeycloakToken) error

// Config stores the configuration to the Keycloak server.
type Config struct {
	URL   string // URL is the token URL of the server
	Realm string // Realm is the realm to use.

	timeout            time.Duration
	fullCertsPath      string
	customClaimsMapper ClaimMapperFunc
	cache              *cache[string, *keyEntry]
}

// TokenContainer stores all relevant token information.
type TokenContainer struct {
	Token         *oauth2.Token
	KeycloakToken *KeycloakToken
}

// NewConfig returns a new configuration for a given Keycloak server.
func NewConfig(url, realm string) *Config {
	return &Config{
		URL:   url,
		Realm: realm,
	}
}

// ConfigOption is an option to configure the Gin middleware.
type ConfigOption func(*Config)

// Timeout is an option how long should be waited as max till
// a fetch call from the Keycloak server should timeout.
func Timeout(timeout time.Duration) ConfigOption {
	return func(cfg *Config) {
		cfg.timeout = timeout
	}
}

// CustomClaimsMapper is an option to apply a custom mapper to the
// incoming token.
func CustomClaimsMapper(fn ClaimMapperFunc) ConfigOption {
	return func(cfg *Config) {
		cfg.customClaimsMapper = fn
	}
}

// FullCertsPath is an option to configure a full path for fetching
// the certificates from the Keycloak server.
func FullCertsPath(path string) ConfigOption {
	return func(cfg *Config) {
		cfg.fullCertsPath = path
	}
}

// Cache is an option how long a fetched certificate should be
// assumed to stay valid.
func Cache(expiration time.Duration) ConfigOption {
	return func(cfg *Config) {
		cfg.cache = newCache[string, *keyEntry](expiration)
	}
}

// With applies all given options to the configuration.
func (cfg *Config) With(options ...ConfigOption) *Config {
	for _, option := range options {
		option(cfg)
	}
	return cfg
}

// Auth returns a Gin middleware given ab access checking function and a configuration.
func Auth(acfn AccessCheckFunction, cfg *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tc, ok := getTokenContainer(ctx, cfg)
		if !ok {
			log.Println("no token in context")
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("no token in context"))
			return
		}
		if !tc.Valid() {
			log.Println("invalid token")
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
			return
		}
		if !acfn(tc, ctx) {
			log.Println("forbidden")
			ctx.AbortWithError(http.StatusForbidden, errors.New("access to resource forbidden"))
			return
		}
		tc.addToContext(ctx)
		log.Println("access granted")
	}
}

func getTokenContainer(ctx *gin.Context, cfg *Config) (*TokenContainer, bool) {
	var (
		oauthToken *oauth2.Token
		tc         *TokenContainer
		err        error
	)

	if oauthToken, err = extractToken(ctx.Request); err != nil {
		slog.Error("[Gin-OAuth] Can not extract oauth2.Token", "err", err)
		return nil, false
	}

	if !oauthToken.Valid() {
		slog.Error("[Gin-OAuth] Invalid Token - nil or expired")
		return nil, false
	}

	if tc, err = buildTokenContainer(oauthToken, cfg); err != nil {
		slog.Error("[Gin-OAuth] Can not extract TokenContainer", "err", err)
		return nil, false
	}

	if tc.KeycloakToken.isExpired() {
		slog.Error("Token expired")
		return nil, false
	}

	return tc, true
}

func extractToken(r *http.Request) (*oauth2.Token, error) {
	hdr := r.Header.Get("Authorization")
	if hdr == "" {
		return nil, errors.New("no authorization header")
	}
	typ, token, ok := strings.Cut(hdr, " ")
	if !ok {
		return nil, errors.New("incomplete authorization header")
	}
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   typ,
	}, nil
}

func buildTokenContainer(token *oauth2.Token, cfg *Config) (*TokenContainer, error) {
	kct, err := decodeToken(token, cfg)
	if err != nil {
		return nil, err
	}
	return &TokenContainer{
		Token: &oauth2.Token{
			AccessToken: token.AccessToken,
			TokenType:   token.TokenType,
		},
		KeycloakToken: kct,
	}, nil
}

func decodeToken(token *oauth2.Token, cfg *Config) (*KeycloakToken, error) {
	kct := KeycloakToken{}

	var err error
	parsedJWT, err := jwt.ParseSigned(token.AccessToken)
	if err != nil {
		slog.Warn("JWT not decodable", "err", err)
		return nil, err
	}
	key, err := getPublicKey(parsedJWT.Headers[0].KeyID, cfg)
	if err != nil {
		slog.Warn("Failed to get publickey", "err", err)
		return nil, err
	}

	if err = parsedJWT.Claims(key, &kct); err != nil {
		slog.Warn("Failed to get claims JWT", "err", err)
		return nil, err
	}

	if cfg.customClaimsMapper != nil {
		claims := func(dst any) error { return parsedJWT.Claims(key, dst) }
		if err = cfg.customClaimsMapper(claims, &kct); err != nil {
			slog.Warn("Failed to get custom claims JWT", "err", err)
			return nil, err
		}
	}

	return &kct, nil
}

func getPublicKey(keyID string, cfg *Config) (any, error) {
	ke, err := fetchPublicKey(keyID, cfg)
	if err != nil {
		return nil, err
	}
	switch kty := strings.ToUpper(ke.Kty); kty {
	case "RSA":
		n, err := base64.RawURLEncoding.DecodeString(ke.N)
		if err != nil {
			return nil, err
		}
		bigN := new(big.Int)
		bigN.SetBytes(n)
		e, err := base64.RawURLEncoding.DecodeString(ke.E)
		if err != nil {
			return nil, err
		}
		bigE := new(big.Int)
		bigE.SetBytes(e)
		return &rsa.PublicKey{
			N: bigN,
			E: int(bigE.Int64()),
		}, nil
	case "EC":
		x, err := base64.RawURLEncoding.DecodeString(ke.X)
		if err != nil {
			return nil, err
		}
		bigX := new(big.Int)
		bigX.SetBytes(x)
		y, err := base64.RawURLEncoding.DecodeString(ke.Y)
		if err != nil {
			return nil, err
		}
		bigY := new(big.Int)
		bigY.SetBytes(y)

		var curve elliptic.Curve
		switch crv := strings.ToUpper(ke.Crv); crv {
		case "P-224":
			curve = elliptic.P224()
		case "P-256":
			curve = elliptic.P256()
		case "P-384":
			curve = elliptic.P384()
		case "P-521":
			curve = elliptic.P521()
		default:
			return nil, fmt.Errorf("EC curve algorithm not supported %q", ke.Kty)
		}
		return &ecdsa.PublicKey{
			Curve: curve,
			X:     bigX,
			Y:     bigY,
		}, nil
	default:
		return nil, fmt.Errorf("no support for keys of type %q", kty)
	}
}

func fetchPublicKey(keyID string, cfg *Config) (*keyEntry, error) {
	if cfg.cache != nil {
		if entry, exists := cfg.cache.get(keyID); exists {
			return entry, nil
		}
	}

	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	if cfg.fullCertsPath != "" {
		u.Path = cfg.fullCertsPath
	} else {
		u.Path = path.Join(u.Path, "realms", cfg.Realm, "protocol/openid-connect/certs")
	}

	client := http.Client{}
	if cfg.timeout != 0 {
		client.Timeout = cfg.timeout
	}

	slog.Debug("requesting keyclock's public key", "url", u)
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot GET public key: %s (%d)",
			resp.Status, resp.StatusCode)
	}

	var certs certs
	if err := func() error {
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(&certs)
	}(); err != nil {
		return nil, err
	}

	for _, entry := range certs.Keys {
		if entry.Kid == keyID {
			if cfg.cache != nil {
				cfg.cache.set(keyID, entry)
			}
			return entry, nil
		}
	}

	return nil, fmt.Errorf("no public key found for kid %q", keyID)
}

// Valid returns true if the given token container is valid.
func (tc *TokenContainer) Valid() bool {
	return tc != nil && tc.Token != nil && tc.Token.Valid()
}

func (tc *TokenContainer) addToContext(ctx *gin.Context) {
	ctx.Set("token", tc.KeycloakToken)
	ctx.Set("uid", tc.KeycloakToken.PreferredUsername)
}
