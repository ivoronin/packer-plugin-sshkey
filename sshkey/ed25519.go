package sshkey

import (
	"crypto"
	"crypto/ed25519"
	"encoding/pem"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
)

const ED25519_BLOCK_TYPE = "OPENSSH PRIVATE KEY"

type ED25519Key struct {
	SSHKey
	pubkey crypto.PublicKey
	key ed25519.PrivateKey
}

func (k *ED25519Key) Generate() error {
	var err error
	k.pubkey, k.key, err = ed25519.GenerateKey(nil)
	if err == nil {
		return err
	}

	return nil
}

func (k *ED25519Key) FromPEM(bytes []byte) error {
	var err error

	key, err := ssh.ParseRawPrivateKey(bytes)
	if err != nil { return err }

	k.key = *key.(*ed25519.PrivateKey)
	k.pubkey = k.key.Public()
	return nil
}

func (k *ED25519Key) ToPEM() ([]byte, error) {
	block := pem.Block{
		Type:  ED25519_BLOCK_TYPE,
		Bytes: edkey.MarshalED25519PrivateKey(k.key),
	}
	pem := pem.EncodeToMemory(&block)
	return pem, nil
}

func (k *ED25519Key) Public() (string, error) {
	return publicKeyStringFor(k.pubkey)
}
