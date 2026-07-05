// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"errors"
	"testing"
)

func TestConfigureAppliesDefaults(t *testing.T) {
	t.Parallel()

	d := Datasource{config: Config{}}

	err := d.Configure(nil)
	if err != nil {
		t.Fatalf("Failed to configure datasource")
	}

	if d.config.Name != "packer" {
		t.Fatalf("Name default = %q, want %q", d.config.Name, "packer")
	}

	if d.config.Type != keyTypeRSA {
		t.Fatalf("Type default = %q, want %q", d.config.Type, keyTypeRSA)
	}

	if d.config.Cache != cacheReuse {
		t.Fatalf("Cache default = %q, want %q", d.config.Cache, cacheReuse)
	}
}

func TestConfigureRejectsUnsupportedCache(t *testing.T) {
	t.Parallel()

	d := Datasource{config: Config{Cache: "fresh"}}

	err := d.Configure(nil)
	if !errors.Is(err, errUnsupportedCache) {
		t.Fatalf("Configure error = %v, want %v", err, errUnsupportedCache)
	}
}
