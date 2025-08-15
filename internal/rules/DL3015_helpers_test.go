// file: internal/rules/DL3015_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// TestRunTokens covers JSON and shell forms.
func TestRunTokens(t *testing.T) {
	if tokens := runTokens(nil); tokens != nil {
		t.Fatalf("expected nil tokens for nil node")
	}
	n := &parser.Node{Value: "run"}
	if tokens := runTokens(n); tokens != nil {
		t.Fatalf("expected nil tokens for missing next")
	}
	jsonNode := &parser.Node{Attributes: map[string]bool{"json": true}, Next: &parser.Node{Value: "echo", Next: &parser.Node{Value: "hi"}}}
	tks := runTokens(jsonNode)
	if len(tks) != 2 || tks[0] != "echo" || tks[1] != "hi" {
		t.Fatalf("unexpected tokens: %#v", tks)
	}
	shellNode := &parser.Node{Next: &parser.Node{Value: "echo hi"}}
	tks = runTokens(shellNode)
	if len(tks) != 2 || tks[0] != "echo" || tks[1] != "hi" {
		t.Fatalf("unexpected shell tokens: %#v", tks)
	}
}

// TestSplitTokens ensures commands are divided at connectors.
func TestSplitTokens(t *testing.T) {
	tokens := []string{"apt-get", "update", "&&", "apt-get", "install", "-y"}
	cmds := splitTokens(tokens)
	if len(cmds) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(cmds))
	}
	if len(splitTokens([]string{"echo"})) != 1 {
		t.Fatalf("expected 1 command")
	}
}

// TestAptInstallMissingFlag covers flag detection.
func TestAptInstallMissingFlag(t *testing.T) {
	if aptInstallMissingFlag([]string{}) {
		t.Fatalf("empty tokens should not flag")
	}
	if aptInstallMissingFlag([]string{"echo"}) {
		t.Fatalf("non-apt command should not flag")
	}
	if !aptInstallMissingFlag([]string{"apt-get", "install"}) {
		t.Fatalf("missing flag should be reported")
	}
	if aptInstallMissingFlag([]string{"apt-get", "install", "--no-install-recommends"}) {
		t.Fatalf("flag present should not report")
	}
	if aptInstallMissingFlag([]string{"apt-get", "install", "-o", "APT::Install-Recommends=false"}) {
		t.Fatalf("option flag should not report")
	}
}
