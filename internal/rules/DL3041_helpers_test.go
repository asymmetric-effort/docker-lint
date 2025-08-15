// file: internal/rules/DL3041_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"
	"testing"
)

// TestIsDnfUpgrade verifies detection of disallowed upgrade commands.
func TestIsDnfUpgrade(t *testing.T) {
	cases := map[string]bool{
		"dnf upgrade":      true,
		"dnf update":       true,
		"dnf -y upgrade":   true,
		"microdnf upgrade": true,
		"dnf install":      false,
		"echo hi":          false,
		"dnf -y":           false,
	}
	for cmd, expect := range cases {
		tokens := strings.Split(cmd, " ")
		if got := isDnfUpgrade(tokens); got != expect {
			t.Fatalf("%s: expected %v got %v", cmd, expect, got)
		}
	}
}
