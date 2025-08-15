// file: internal/rules/DL3014_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"
	"testing"
)

// TestIsAptGetInstall checks apt-get install detection.
func TestIsAptGetInstall(t *testing.T) {
	cases := map[string]bool{
		"apt-get install":    true,
		"apt-get -y install": true,
		"apt-get update":     false,
		"echo":               false,
	}
	for cmd, expect := range cases {
		tokens := strings.Split(cmd, " ")
		if got := isAptGetInstall(tokens); got != expect {
			t.Fatalf("%s: expected %v got %v", cmd, expect, got)
		}
	}
}

// TestHasYesOption verifies detection of non-interactive flags.
func TestHasYesOption(t *testing.T) {
	cases := map[string]bool{
		"apt-get install -y":              true,
		"apt-get install --yes":           true,
		"apt-get install --assume-yes":    true,
		"apt-get install -qq":             true,
		"apt-get install -q=2":            true,
		"apt-get install --quiet=2":       true,
		"apt-get install -q -q":           true,
		"apt-get install --quiet --quiet": true,
		"apt-get install":                 false,
	}
	for cmd, expect := range cases {
		tokens := strings.Split(cmd, " ")
		if got := hasYesOption(tokens); got != expect {
			t.Fatalf("%s: expected %v got %v", cmd, expect, got)
		}
	}
}
