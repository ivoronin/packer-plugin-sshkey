// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const testKeyName = "test"

// Execute writes the key under PACKER_CACHE_DIR on first run and reads it back
// on the second. This test mutates the process environment via t.Setenv, so it
// cannot run in parallel.
func TestExecuteGeneratesThenLoads(t *testing.T) {
	t.Setenv("PACKER_CACHE_DIR", t.TempDir())

	d := Datasource{config: Config{Name: testKeyName, Type: keyTypeED25519, Cache: cacheReuse}}

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

func TestExecuteRejectsUnconfiguredCache(t *testing.T) {
	t.Setenv("PACKER_CACHE_DIR", t.TempDir())

	d := Datasource{config: Config{Name: testKeyName, Type: keyTypeED25519}}

	_, err := d.Execute()
	if !errors.Is(err, errUnsupportedCache) {
		t.Fatalf("Execute error = %v, want %v", err, errUnsupportedCache)
	}
}

// Execute with cache="unique" writes a fresh key under PACKER_CACHE_DIR on
// every run. This test mutates the process environment via t.Setenv, so it
// cannot run in parallel.
func TestExecuteUniqueCacheCreatesNewKeyEveryRun(t *testing.T) {
	cacheDir := t.TempDir()
	t.Setenv("PACKER_CACHE_DIR", cacheDir)

	d := Datasource{config: Config{Name: testKeyName, Type: keyTypeED25519, Cache: cacheUnique}}

	first, err := d.Execute()
	if err != nil {
		t.Fatalf("first Execute: %v", err)
	}

	second, err := d.Execute()
	if err != nil {
		t.Fatalf("second Execute: %v", err)
	}

	firstPath := first.GetAttr("private_key_path").AsString()
	secondPath := second.GetAttr("private_key_path").AsString()

	if firstPath == secondPath {
		t.Fatalf("cache path was reused in unique mode: %q", firstPath)
	}

	if filepath.Dir(firstPath) != cacheDir {
		t.Fatalf("first private key dir = %q, want %q", filepath.Dir(firstPath), cacheDir)
	}

	if filepath.Dir(secondPath) != cacheDir {
		t.Fatalf("second private key dir = %q, want %q", filepath.Dir(secondPath), cacheDir)
	}

	if _, err := os.Stat(firstPath); err != nil {
		t.Fatalf("stat first private key: %v", err)
	}

	if _, err := os.Stat(secondPath); err != nil {
		t.Fatalf("stat second private key: %v", err)
	}

	if first.GetAttr("public_key").AsString() == second.GetAttr("public_key").AsString() {
		t.Fatal("public_key was reused in unique mode")
	}
}
