package sshkey

import (
	_ "embed"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"os/exec"
	"testing"
)

//go:embed test-fixtures/rsa.pkr.hcl
var testDatasourceTemplateRSA string

//go:embed test-fixtures/ed25519.pkr.hcl
var testDatasourceTemplateED25519 string

func check(buildCommand *exec.Cmd, logfile string) error {
	if buildCommand.ProcessState != nil {
		if buildCommand.ProcessState.ExitCode() != 0 {
			return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
		}
	}
	return nil
}

func TestAccSSHKeyManager(t *testing.T) {
	testGenerateRSA := &acctest.PluginTestCase{
		Name:     "sshkey_datasource_create_rsa",
		Template: testDatasourceTemplateRSA,
		Check:    check,
	}
	testLoadRSA := &acctest.PluginTestCase{
		Name:     "sshkey_datasource_load_rsa",
		Template: testDatasourceTemplateRSA,
		Check:    check,
	}
	testGenerateED25519 := &acctest.PluginTestCase{
		Name:     "sshkey_datasource_create_ed25519",
		Template: testDatasourceTemplateED25519,
		Check:    check,
	}
	testLoadED25519 := &acctest.PluginTestCase{
		Name:     "sshkey_datasource_load_ed25519",
		Template: testDatasourceTemplateED25519,
		Check:    check,
	}
	acctest.TestPlugin(t, testGenerateRSA)
	acctest.TestPlugin(t, testLoadRSA)
	acctest.TestPlugin(t, testGenerateED25519)
	acctest.TestPlugin(t, testLoadED25519)
}
