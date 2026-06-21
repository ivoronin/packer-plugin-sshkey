// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
)

const ed25519BlockType = "PRIVATE KEY"

var errNotED25519 = errors.New("not an ed25519 private key")

var _ SSHKey = (*ED25519Key)(nil)

// ED25519Key wraps an Ed25519 private key and implements the SSHKey interface.
type ED25519Key struct {
	key ed25519.PrivateKey
}

// Generate creates a new Ed25519 key pair using the system CSPRNG.
func (k *ED25519Key) Generate() error {
	return k.generate(rand.Reader)
}

// FromPEM parses a PEM-encoded PKCS#8 private key and asserts it is Ed25519.
func (k *ED25519Key) FromPEM(bytes []byte) error {
	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != ed25519BlockType {
		return errDecodePEM
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	key, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return errNotED25519
	}

	k.key = key

	return nil
}

// ToPEM serialises the Ed25519 private key as a PEM-encoded PKCS#8 block.
func (k *ED25519Key) ToPEM() ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(k.key)
	if err != nil {
		return nil, err
	}

	block := pem.Block{
		Type:  ed25519BlockType,
		Bytes: der,
	}

	return pem.EncodeToMemory(&block), nil
}

// Public returns the OpenSSH-format authorised-keys string for the Ed25519 key.
func (k *ED25519Key) Public() (string, error) {
	return publicKeyStringFor(k.key.Public())
}

func (k *ED25519Key) generate(random io.Reader) error {
	_, key, err := ed25519.GenerateKey(random)
	if err != nil {
		return err
	}

	k.key = key

	return nil
}
