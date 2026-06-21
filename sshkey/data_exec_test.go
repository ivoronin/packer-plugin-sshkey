// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"testing"
)

// Execute writes the key under PACKER_CACHE_DIR on first run and reads it back
// on the second. This test mutates the process environment via t.Setenv, so it
// cannot run in parallel.
func TestExecuteGeneratesThenLoads(t *testing.T) {
	t.Setenv("PACKER_CACHE_DIR", t.TempDir())

	d := Datasource{config: Config{Name: "test", Type: "ed25519"}}

	first, err := d.Execute()
	if err != nil {
		t.Fatalf("first Execute: %v", err)
	}

	second, err := d.Execute()
	if err != nil {
		t.Fatalf("second Execute (load from cache): %v", err)
	}

	firstPath := first.GetAttr("private_key_path").AsString()

	secondPath := second.GetAttr("private_key_path").AsString()

	if firstPath != secondPath {
		t.Fatalf("cache path changed between runs: %q vs %q", firstPath, secondPath)
	}

	pub := second.GetAttr("public_key").AsString()
	if pub == "" {
		t.Fatal("public_key is empty")
	}

	if first.GetAttr("public_key").AsString() != pub {
		t.Fatal("public_key changed between cached runs; second Execute regenerated instead of loading from cache")
	}
}
