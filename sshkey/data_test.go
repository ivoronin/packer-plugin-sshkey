// Copyright (c) Ilya Voronin
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"testing"
)

func TestDatasource(t *testing.T) {
	t.Parallel()

	d := Datasource{config: Config{}}

	err := d.Configure(nil)
	if err != nil {
		t.Fatalf("Failed to configure datasource")
	}
}
