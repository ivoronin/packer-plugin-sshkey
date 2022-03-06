//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package sshkey

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Name string `mapstructure:"name"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKey      string `mapstructure:"public_key"`
}

func (d *Datasource) Configure(raws ...interface{}) error {
	if err := config.Decode(&d.config, nil, raws...); err != nil {
		return err
	}

	return nil
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	var privateKey ed25519.PrivateKey
	var err error

	keyName := d.config.Name
	if keyName == "" {
		keyName = "packer"
	}

	privateKeyNameSuffix := strings.ReplaceAll(keyName, string(os.PathSeparator), "_")
	privateKeyName := "ssh_private_key_" + privateKeyNameSuffix + ".pem"

	privateKeyPath, err := packer.CachePath(privateKeyName)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	privateKeyPEM, err := ioutil.ReadFile(privateKeyPath)
	if err == nil {
		if privateKey, err = decodePrivateKeyFromPEM(privateKeyPEM); err != nil {
			return cty.NullVal(cty.EmptyObject), err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if privateKey, err = generatePrivateKey(); err != nil {
			return cty.NullVal(cty.EmptyObject), err
		}

		privateKeyPEM, err = encodePrivateKeyToPEM(privateKey)
		if err != nil {
			return cty.NullVal(cty.EmptyObject), err
		}
		if err = ioutil.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
			return cty.NullVal(cty.EmptyObject), err
		}
	} else {
		return cty.NullVal(cty.EmptyObject), err
	}

	publicKeyString, err := generatePublicKeyString(privateKey.Public())
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	output := DatasourceOutput{
		PrivateKeyPath: privateKeyPath,
		PublicKey:      publicKeyString + " " + keyName,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}

func generatePrivateKey() (ed25519.PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey ed25519.PrivateKey) ([]byte, error) {
	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privBlock := pem.Block{
		Type:    "PRIVATE KEY",
		Bytes:   privDER,
	}

	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM, nil
}

func decodePrivateKeyFromPEM(privateKeyPEM []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("unable to decode PEM")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(ed25519.PrivateKey), nil
}

func generatePublicKeyString(publicKey crypto.PublicKey) (string, error) {
	rsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	keyBytes := ssh.MarshalAuthorizedKey(rsaKey)
	keyString := strings.TrimRight(string(keyBytes[:]), "\r\n")

	return keyString, nil
}
