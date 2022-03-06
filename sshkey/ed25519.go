package sshkey

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const ED25519_BLOCK_TYPE = "PRIVATE KEY"

type ED25519Key struct {
	SSHKey
	key ed25519.PrivateKey
}

func (k *ED25519Key) Generate() error {
	var err error
	_, k.key, err = ed25519.GenerateKey(nil)
	if err == nil {
		return err
	}

	return nil
}

func (k *ED25519Key) FromPEM(bytes []byte) error {
	var err error
	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != ED25519_BLOCK_TYPE {
		return errors.New("unable to decode PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	k.key = key.(ed25519.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}

func (k *ED25519Key) ToPEM() ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(k.key)
	if err != nil {
		return nil, err
	}
	block := pem.Block{
		Type:  ED25519_BLOCK_TYPE,
		Bytes: der,
	}
	pem := pem.EncodeToMemory(&block)
	return pem, nil
}

func (k *ED25519Key) Public() (string, error) {
	return publicKeyStringFor(k.key.Public())
}
