// file: internal/rules/DL3043_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import "testing"

// TestNeedsOSVersionTag exercises edge cases for tag requirements.
func TestNeedsOSVersionTag(t *testing.T) {
	cases := map[string]bool{
		"ubuntu":             true,
		"ubuntu:latest":      true,
		"ubuntu:jammy":       true,
		"ubuntu:22.04":       false,
		"alpine@sha256:dead": true,
		"$BASE":              false,
		"scratch":            false,
	}
	for image, expect := range cases {
		if got := needsOSVersionTag(image); got != expect {
			t.Fatalf("%s: expected %v got %v", image, expect, got)
		}
	}
}

// TestIsOSImage verifies detection of OS base images.
func TestIsOSImage(t *testing.T) {
	cases := map[string]bool{
		"ubuntu":         true,
		"library/alpine": true,
		"ghcr.io/other":  false,
		"golang":         false,
	}
	for image, expect := range cases {
		if got := isOSImage(image); got != expect {
			t.Fatalf("%s: expected %v got %v", image, expect, got)
		}
	}
}
