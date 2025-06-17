// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ISDuBA/ISDuBA/internal/config"
)

const writeKeyMsg = `================================================================
No [sources] aes_key in configuration found. Using newly generated:
%s
Store this under
  [sources]
  aes_key = "<hex string above>"
in your configuration or in a separate file referencing it with
  [sources]
  aes_key = "@/path/to/file"
or the encrypted database fields will not be accessible any more 
on next boot. See documentation for details.
================================================================
`

func createCipherKey(cfg *config.Config) ([]byte, error) {
	aesKey := cfg.Sources.AESKey
	var err error
	var key []byte
	switch {
	case aesKey == "":
		// No key given -> Create new one and write to STDOUT.
		key = make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("cannot read random key: %w", err)
		}
		fmt.Printf(writeKeyMsg, hex.EncodeToString(key))
	case aesKey[0] == '@':
		fname := aesKey[1:]
		if key, err = loadKeyFromFile(fname); err != nil {
			return nil, fmt.Errorf("loading key from file %q failed: %w", fname, err)
		}
	default:
		if key, err = hex.DecodeString(aesKey); err != nil {
			return nil, fmt.Errorf("sources.aes_key is invalid: %w", err)
		}
	}
	if len(key) < 32 {
		return nil, errors.New("key is too short")
	}
	return key[:32], nil
}

func loadKeyFromFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var key []byte
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		if key, err = hex.DecodeString(line); err != nil {
			return nil, err
		}
		break
	}
	return key, sc.Err()
}

func (m *Manager) createCipher() (cipher.AEAD, error) {
	blockCipher, err := aes.NewCipher(m.cipherKey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(blockCipher)
}

// encrypt encrypts data with the key of the manager.
func (m *Manager) encrypt(data []byte) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	cipher, err := m.createCipher()
	if err != nil {
		return nil, err
	}
	cipherLen := len(data) + 16 - len(data)%16
	nonce := make([]byte, cipher.NonceSize(), cipher.NonceSize()+cipherLen)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("reading nonce failed: %w", err)
	}
	return cipher.Seal(nonce, nonce, data, nil), nil
}

// decrypt decrypts data with the key of the manager.
func (m *Manager) decrypt(data []byte) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	cipher, err := m.createCipher()
	if err != nil {
		return nil, err
	}
	nonce, cipherText := data[:cipher.NonceSize()], data[cipher.NonceSize():]
	return cipher.Open(nil, nonce, cipherText, nil)
}
