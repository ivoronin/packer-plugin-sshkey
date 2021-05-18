package sshkey

import (
	"testing"
)

func TestDatasource(t *testing.T) {
	d := Datasource{config: Config{}}
	err := d.Configure(nil)
	if err != nil {
		t.Fatalf("Failed to configure datasource")
	}
}
