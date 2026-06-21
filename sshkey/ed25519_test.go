// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

func TestED25519RoundTrip(t *testing.T) {
	t.Parallel()

	var k ED25519Key
	if err := k.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	pem, err := k.ToPEM()
	if err != nil {
		t.Fatalf("ToPEM: %v", err)
	}

	var loaded ED25519Key
	if err := loaded.FromPEM(pem); err != nil {
		t.Fatalf("FromPEM: %v", err)
	}

	got, err := loaded.Public()
	if err != nil {
		t.Fatalf("Public: %v", err)
	}

	if !strings.HasPrefix(got, "ssh-ed25519 ") {
		t.Fatalf("public key does not start with ssh-ed25519: %q", got)
	}
}

func TestED25519GeneratePropagatesError(t *testing.T) {
	t.Parallel()

	var k ED25519Key

	err := k.generate(iotest.ErrReader(errors.New("boom")))
	if err == nil {
		t.Fatal("expected Generate to propagate the rand reader error, got nil")
	}
}

func TestED25519FromPEMRejectsWrongKeyType(t *testing.T) {
	t.Parallel()

	var rk RSAKey
	if err := rk.Generate(); err != nil {
		t.Fatalf("rsa Generate: %v", err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(rk.key)
	if err != nil {
		t.Fatalf("MarshalPKCS8PrivateKey: %v", err)
	}

	wrongPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	var ed ED25519Key
	if err := ed.FromPEM(wrongPEM); !errors.Is(err, errNotED25519) {
		t.Fatalf("expected errNotED25519, got %v", err)
	}
}
