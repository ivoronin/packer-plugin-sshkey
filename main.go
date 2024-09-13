package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
	"github.com/ivoronin/packer-plugin-sshkey/sshkey"
)

var (
	Version           = "1.2.2"
	VersionPrerelease = ""
	PluginVersion     = version.NewPluginVersion(Version, VersionPrerelease, "")
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(sshkey.Datasource))
	pps.SetVersion(PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
