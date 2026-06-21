// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package sshkey

import (
	"errors"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

var errUnsupportedKeyType = errors.New("unsupported key type")

// Config holds the user-supplied HCL configuration for the sshkey data source.
type Config struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
}

// Datasource is the Packer data source that generates or loads an SSH key pair.
type Datasource struct {
	config Config
}

// DatasourceOutput is the value the data source exposes to downstream Packer config.
type DatasourceOutput struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKey      string `mapstructure:"public_key"`
}

// Configure decodes HCL config into the datasource and applies defaults.
func (d *Datasource) Configure(raws ...any) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
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

// ConfigSpec returns the HCL object spec for the data source inputs.
func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

// OutputSpec returns the HCL object spec for the data source outputs.
func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

// Execute generates a new key or loads an existing one from the Packer cache.
func (d *Datasource) Execute() (cty.Value, error) {
	nullVal := cty.NullVal(cty.EmptyObject)

	key, err := newKey(d.config.Type)
	if err != nil {
		return nullVal, err
	}

	keyTag := strings.ReplaceAll(d.config.Name, string(os.PathSeparator), "_")
	keyName := "ssh_private_key_" + keyTag + "_" + d.config.Type + ".pem"

	keyPath, err := packer.CachePath(keyName)
	if err != nil {
		return nullVal, err
	}

	err = loadOrGenerate(key, keyPath)
	if err != nil {
		return nullVal, err
	}

	pubKeyStr, err := key.Public()
	if err != nil {
		return nullVal, err
	}

	output := DatasourceOutput{
		PrivateKeyPath: keyPath,
		PublicKey:      pubKeyStr + " " + keyTag,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}

//nolint:ireturn // factory must return the SSHKey interface to select between RSAKey and ED25519Key
func newKey(keyType string) (SSHKey, error) {
	switch keyType {
	case "rsa":
		return new(RSAKey), nil
	case "ed25519":
		return new(ED25519Key), nil
	default:
		return nil, errUnsupportedKeyType
	}
}

func loadOrGenerate(key SSHKey, keyPath string) error {
	pemBytes, err := os.ReadFile(keyPath) //nolint:gosec // G304: keyPath comes from packer.CachePath

	switch {
	case err == nil:
		return key.FromPEM(pemBytes)
	case errors.Is(err, os.ErrNotExist):
		return generateAndSave(key, keyPath)
	default:
		return err
	}
}

func generateAndSave(key SSHKey, keyPath string) error {
	err := key.Generate()
	if err != nil {
		return err
	}

	pemBytes, err := key.ToPEM()
	if err != nil {
		return err
	}

	return os.WriteFile(keyPath, pemBytes, 0o600)
}
