// file: internal/rules/DL3038_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"
	"testing"
)

// TestIsDnfInstall verifies detection of install subcommands.
func TestIsDnfInstall(t *testing.T) {
	cases := map[string]bool{
		"dnf install":           true,
		"dnf groupinstall":      true,
		"microdnf localinstall": true,
		"dnf update":            false,
		"echo hi":               false,
		"dnf":                   false,
	}
	for cmd, expect := range cases {
		tokens := strings.Split(cmd, " ")
		if got := isDnfInstall(tokens); got != expect {
			t.Fatalf("%s: expected %v got %v", cmd, expect, got)
		}
	}
}
