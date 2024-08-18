// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func (s *source) updateCertificate() error {
	private, err := s.clientCertPrivate()
	if err != nil {
		s.tlsCertificates = nil
		return err
	}
	if s.clientCertPublic == nil || private == nil {
		s.tlsCertificates = nil
		return nil
	}
	passphrase, err := s.clientCertPassphrase()
	if err != nil {
		s.tlsCertificates = nil
		return err
	}
	if passphrase != nil {
		block, _ := pem.Decode(private)
		if block == nil {
			s.tlsCertificates = nil
			return errors.New("private key has no PEM block")
		}
		//lint:ignore SA1019 This is insecure by design.
		keyDER, err := x509.DecryptPEMBlock(block, passphrase)
		if err != nil {
			s.tlsCertificates = nil
			return err
		}
		// Update keyBlock with the plaintext bytes and clear the now obsolete
		// headers.
		block.Bytes = keyDER
		block.Headers = nil

		// Turn the key back into PEM format so we can leverage tls.X509KeyPair,
		// which will deal with the intricacies of error handling, different key
		// types, certificate chains, etc
		private = pem.EncodeToMemory(block)
	}
	cert, err := tls.X509KeyPair(s.clientCertPublic, private)
	if err != nil {
		s.tlsCertificates = nil
		return err
	}
	s.tlsCertificates = []tls.Certificate{cert}
	return nil
}
