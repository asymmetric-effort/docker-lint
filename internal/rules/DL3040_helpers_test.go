// file: internal/rules/DL3040_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import "testing"

// TestCleansDnfCache exercises cache cleanup detection.
func TestCleansDnfCache(t *testing.T) {
	if !cleansDnfCache([]string{"dnf", "clean", "all"}) {
		t.Fatalf("expected dnf clean all to be detected")
	}
	if !cleansDnfCache([]string{"rm", "-fr", "/var/cache/dnf"}) {
		t.Fatalf("expected rm -fr to be detected")
	}
	if cleansDnfCache([]string{"rm", "/var/cache/dnf"}) {
		t.Fatalf("missing flags should not detect")
	}
	if !cleansDnfCache([]string{"find", "/var/cache/dnf", "-delete"}) {
		t.Fatalf("expected find -delete detection")
	}
	if cleansDnfCache([]string{"find", "/var/cache/dnf", "-print"}) {
		t.Fatalf("-print should not detect")
	}
	if cleansDnfCache([]string{"echo", "hi"}) {
		t.Fatalf("unrelated command should not detect")
	}
}
