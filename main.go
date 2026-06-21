// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

// Command packer-plugin-sshkey is a Packer plugin that generates SSH keys.
package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/ivoronin/packer-plugin-sshkey/sshkey"
	"github.com/ivoronin/packer-plugin-sshkey/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(sshkey.Datasource))
	pps.SetVersion(version.PluginVersion)

	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
