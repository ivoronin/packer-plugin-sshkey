// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
)

const (
	rsaBits      = 2048
	rsaBlockType = "RSA PRIVATE KEY"
)

var errDecodePEM = errors.New("unable to decode PEM")

var _ SSHKey = (*RSAKey)(nil)

// RSAKey wraps an RSA private key and implements the SSHKey interface.
type RSAKey struct {
	key *rsa.PrivateKey
}

// Generate creates a new RSA key pair using the system CSPRNG.
func (k *RSAKey) Generate() error {
	return k.generate(rand.Reader)
}

// FromPEM parses a PEM-encoded PKCS#1 RSA private key.
func (k *RSAKey) FromPEM(bytes []byte) error {
	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != rsaBlockType {
		return errDecodePEM
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	k.key = key

	return nil
}

// ToPEM serialises the RSA private key as a PEM-encoded PKCS#1 block.
func (k *RSAKey) ToPEM() ([]byte, error) {
	der := x509.MarshalPKCS1PrivateKey(k.key)
	block := pem.Block{
		Type:  rsaBlockType,
		Bytes: der,
	}

	return pem.EncodeToMemory(&block), nil
}

// Public returns the OpenSSH-format authorised-keys string for the RSA key.
func (k *RSAKey) Public() (string, error) {
	return publicKeyStringFor(k.key.Public())
}

func (k *RSAKey) generate(random io.Reader) error {
	key, err := rsa.GenerateKey(random, rsaBits)
	if err != nil {
		return err
	}

	k.key = key

	return nil
}
