// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

// Package sshkey implements a Packer data source that generates and caches SSH keys.
package sshkey

import (
	"crypto"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SSHKey is the common interface implemented by RSAKey and ED25519Key.
type SSHKey interface {
	Generate() error
	ToPEM() ([]byte, error)
	FromPEM(bytes []byte) error
	Public() (string, error)
}

func publicKeyStringFor(privKey crypto.PrivateKey) (string, error) {
	pubKey, err := ssh.NewPublicKey(privKey)
	if err != nil {
		return "", err
	}

	bytes := ssh.MarshalAuthorizedKey(pubKey)
	str := strings.TrimRight(string(bytes), "\r\n")

	return str, nil
}
