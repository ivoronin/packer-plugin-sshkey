package sshkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const RSA_BITS = 2048
const RSA_BLOCK_TYPE = "RSA PRIVATE KEY"

type RSAKey struct {
	SSHKey
	key *rsa.PrivateKey
}

func (k *RSAKey) Generate() error {
	var err error
	k.key, err = rsa.GenerateKey(rand.Reader, RSA_BITS)
	if err == nil {
		return err
	}

	return nil
}

func (k *RSAKey) FromPEM(bytes []byte) error {
	var err error
	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != RSA_BLOCK_TYPE {
		return errors.New("unable to decode PEM")
	}

	k.key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	return nil
}

func (k *RSAKey) ToPEM() ([]byte, error) {
	der := x509.MarshalPKCS1PrivateKey(k.key)
	block := pem.Block{
		Type:  RSA_BLOCK_TYPE,
		Bytes: der,
	}
	pem := pem.EncodeToMemory(&block)
	return pem, nil
}

func (k *RSAKey) Public() (string, error) {
	return publicKeyStringFor(k.key.Public())
}
