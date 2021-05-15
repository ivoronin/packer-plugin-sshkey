//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package sshkey

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"golang.org/x/crypto/ssh"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

type Config struct {
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKey  string `mapstructure:"public_key"`
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
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
	var privateKey *rsa.PrivateKey
	var privateKeyPEM []byte

	privateKeyPath, err := packer.CachePath("ssh_private_key.pem")

	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	privateKeyPEM, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		privateKey, err = generatePrivateKey(2048)
		if err != nil {
			return cty.NullVal(cty.EmptyObject), err
		}
		privateKeyPEM = encodePrivateKeyToPEM(privateKey)
		err = ioutil.WriteFile(privateKeyPath, privateKeyPEM, 0600)
	} else {
		privateKey, err = decodePrivateKeyFromPEM(privateKeyPEM)
		if err != nil {
			return cty.NullVal(cty.EmptyObject), err
		}
	}
	publicKeyString, err := generatePublicKeyString(&privateKey.PublicKey)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	output := DatasourceOutput{
		PrivateKeyPath: privateKeyPath,
		PublicKey: publicKeyString,
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func decodePrivateKeyFromPEM(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Unable to decode PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func generatePublicKeyString(publicKey *rsa.PublicKey) (string, error) {
	rsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	keyBytes := ssh.MarshalAuthorizedKey(rsaKey)

	return string(keyBytes[:]), nil
}
