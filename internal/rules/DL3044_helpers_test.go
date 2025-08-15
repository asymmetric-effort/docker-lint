// file: internal/rules/DL3044_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"
	"testing"
)

// TestHasUnpinnedDnfInstall exercises detection of unpinned installs.
func TestHasUnpinnedDnfInstall(t *testing.T) {
	cases := map[string]bool{
		"dnf install pkg-1":      false,
		"dnf install pkg":        true,
		"microdnf install pkg-1": false,
		"dnf --best install pkg": true,
		"dnf pkg install":        false,
		"echo":                   false,
	}
	for cmd, expect := range cases {
		tokens := strings.Split(cmd, " ")
		if got := hasUnpinnedDnfInstall(tokens); got != expect {
			t.Fatalf("%s: expected %v got %v", cmd, expect, got)
		}
	}
}
