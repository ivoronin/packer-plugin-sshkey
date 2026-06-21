// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

// Package version defines the plugin version, overridable at build time via ldflags.
package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	// Version is the main version number.
	Version = "1.3.0"

	// VersionPrerelease is a marker for the version (e.g. "dev", "beta"); empty for releases.
	VersionPrerelease = ""

	// VersionMetadata is optional build metadata.
	VersionMetadata = ""

	// PluginVersion is the assembled plugin version handed to the SDK.
	PluginVersion = version.NewPluginVersion(Version, VersionPrerelease, VersionMetadata)
)
