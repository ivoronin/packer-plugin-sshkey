//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package sshkey

import (
	"errors"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
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

	if d.config.Name == "" {
		d.config.Name = "packer"
	}

	if d.config.Type == "" {
		d.config.Type = "rsa"
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
	var err error
	var key SSHKey

	nv := cty.NullVal(cty.EmptyObject)

	switch d.config.Type {
	case "rsa":
		key = new(RSAKey)
	case "ed25519":
		key = new(ED25519Key)
	default:
		return nv, errors.New("unsupported key type")
	}

	keyTag := strings.ReplaceAll(d.config.Name, string(os.PathSeparator), "_")
	keyName := "ssh_private_key_" + keyTag + "_" + d.config.Type + ".pem"

	keyPath, err := packer.CachePath(keyName)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	pem, err := ioutil.ReadFile(keyPath)
	if err == nil {
		if err = key.FromPEM(pem); err != nil {
			return nv, err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err = key.Generate(); err != nil {
			return nv, err
		}

		pem, err = key.ToPEM()
		if err != nil {
			return nv, err
		}
		if err = ioutil.WriteFile(keyPath, pem, 0600); err != nil {
			return nv, err
		}
	} else {
		return nv, err
	}

	pubKeyStr, err := key.Public()
	if err != nil {
		return nv, err
	}

	output := DatasourceOutput{
		PrivateKeyPath: keyPath,
		PublicKey:      pubKeyStr + " " + keyTag,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
