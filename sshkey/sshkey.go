package sshkey

import (
	"crypto"
	"golang.org/x/crypto/ssh"
	"strings"
)

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
	str := strings.TrimRight(string(bytes[:]), "\r\n")

	return str, nil
}
