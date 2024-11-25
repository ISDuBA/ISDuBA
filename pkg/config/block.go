// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package config

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"syscall"
	"time"
)

// PortRange is a range of ports.
type PortRange [2]int

// IPRange represents a net address range.
type IPRange struct {
	*net.IPNet
}

// UnmarshalTOML implements [toml.Unmarshaler].
func (pr *PortRange) UnmarshalTOML(data any) error {
	switch v := data.(type) {
	case int64:
		(*pr)[0], (*pr)[1] = int(v), int(v)
	case []any:
		if len(v) != 2 {
			return errors.New("invalid length")
		}
		a, ok1 := v[0].(int64)
		b, ok2 := v[1].(int64)
		if !ok1 || !ok2 {
			return errors.New("invalid range type")
		}
		(*pr)[0], (*pr)[1] = int(min(a, b)), int(max(a, b))
	default:
		return fmt.Errorf("unsupported type: %T", data)
	}
	return nil
}

// Contains checks if a given port is in the range.
func (pr PortRange) Contains(port int) bool {
	return pr[0] <= port && port <= pr[1]
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (br *IPRange) UnmarshalText(text []byte) error {
	cidr := string(text)
	_, blocked, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	br.IPNet = blocked
	return nil
}

func (g *General) blockedIP(ip net.IP) bool {
	if g.BlockLoopback && (ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()) {
		return true
	}
	for _, blocked := range g.BlockedRanges {
		if blocked.Contains(ip) {
			for _, allowed := range g.AllowedIPs {
				if allowed.Equal(ip) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (g *General) allowedPort(port int) bool {
	for _, r := range g.AllowedPorts {
		if r.Contains(port) {
			return true
		}
	}
	return false
}

func (g *General) controlDialing(_, address string, _ syscall.RawConn) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	// Check if the port is allowed
	if len(g.AllowedPorts) > 0 {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid port: %q", host)
		}
		if !g.allowedPort(p) {
			return fmt.Errorf("port %d is not an allowed port", p)
		}
	}
	// Check if the IP is blocked.
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("invalid IP: %q", host)
	}
	if g.blockedIP(ip) {
		return errors.New("accessing address is not allowed")
	}
	return nil
}

// Transport returns an [http.DefaultTransport] like [http.Transport] with
// an installed dialing control to limit access to the configured constraints.
func (g *General) Transport() *http.Transport {
	// This mainly an http.DefaultTransport with a dialer control.
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Control:   g.controlDialing,
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
