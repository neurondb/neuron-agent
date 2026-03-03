/*-------------------------------------------------------------------------
 *
 * root_test.go
 *    Tests for CLI root command.
 *
 *-------------------------------------------------------------------------
 */

package cmd

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}
	if rootCmd.Use != "neuronagent-cli" {
		t.Errorf("rootCmd.Use = %q", rootCmd.Use)
	}
	if len(rootCmd.Commands()) == 0 {
		t.Error("rootCmd should have subcommands")
	}
}
