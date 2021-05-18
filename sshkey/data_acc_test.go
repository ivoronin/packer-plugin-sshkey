package sshkey

import (
	_ "embed"
	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"os/exec"
	"fmt"
	"testing"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceTemplate string

func TestAccSSHKeyManager(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "sshkey_datasource_test",
		Template: testDatasourceTemplate,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
