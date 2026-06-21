// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

func TestRSARoundTrip(t *testing.T) {
	t.Parallel()

	var k RSAKey
	if err := k.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	pem, err := k.ToPEM()
	if err != nil {
		t.Fatalf("ToPEM: %v", err)
	}

	var loaded RSAKey
	if err := loaded.FromPEM(pem); err != nil {
		t.Fatalf("FromPEM: %v", err)
	}

	got, err := loaded.Public()
	if err != nil {
		t.Fatalf("Public: %v", err)
	}

	if !strings.HasPrefix(got, "ssh-rsa ") {
		t.Fatalf("public key does not start with ssh-rsa: %q", got)
	}
}

func TestRSAGeneratePropagatesError(t *testing.T) {
	t.Parallel()

	var k RSAKey

	err := k.generate(iotest.ErrReader(errors.New("boom")))
	if err == nil {
		t.Fatal("expected Generate to propagate the rand reader error, got nil")
	}
}
