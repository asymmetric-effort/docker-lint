// file: internal/rules/DL3013_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"testing"
)

// TestSplitRunCommands covers token splitting with connectors.
func TestSplitRunCommands(t *testing.T) {
	tokens := []string{"pip", "install", "a", "&&", "pip", "install", "b", ";", "echo", "done"}
	cmds := splitRunCommands(tokens)
	if len(cmds) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(cmds))
	}
	if len(cmds[0]) != 3 || len(cmds[1]) != 3 || len(cmds[2]) != 2 {
		t.Fatalf("unexpected command segments: %#v", cmds)
	}
	if c := splitRunCommands([]string{"echo", "hi"}); len(c) != 1 {
		t.Fatalf("expected single command, got %d", len(c))
	}
}

// TestPipInstallIndex exercises pipInstallIndex for various forms.
func TestPipInstallIndex(t *testing.T) {
	if idx, ok := pipInstallIndex([]string{"pip", "install", "pkg"}); !ok || idx != 2 {
		t.Fatalf("expected index 2, got %d %v", idx, ok)
	}
	if idx, ok := pipInstallIndex([]string{"python3", "-m", "pip", "install", "pkg"}); !ok || idx != 4 {
		t.Fatalf("expected index 4, got %d %v", idx, ok)
	}
	if _, ok := pipInstallIndex([]string{"pip"}); ok {
		t.Fatalf("expected not ok for incomplete command")
	}
}

// TestPipVersionFixed ensures detection of version pinning.
func TestPipVersionFixed(t *testing.T) {
	cases := map[string]bool{
		"pkg@git+https://repo": true,
		"pkg==1.0":             true,
		"pkg>=1.0":             true,
		"pkg.whl":              true,
		"pkg.tar.gz":           true,
		"path/to/pkg":          true,
		"pkg":                  false,
	}
	for pkg, expect := range cases {
		if got := pipVersionFixed(pkg); got != expect {
			t.Fatalf("pkg %s expected %v got %v", pkg, expect, got)
		}
	}
}

// TestViolatesPipPin exercises edge cases in version checks.
func TestViolatesPipPin(t *testing.T) {
	if violatesPipPin([]string{"pip", "install", "pkg"}) != true {
		t.Fatalf("expected violation")
	}
	if violatesPipPin([]string{"pip", "install", "pkg==1"}) {
		t.Fatalf("unexpected violation")
	}
	if violatesPipPin([]string{"pip", "install", "-r", "req.txt"}) {
		t.Fatalf("requirement should not violate")
	}
	if violatesPipPin([]string{"pip", "install", "--constraint", "c.txt"}) {
		t.Fatalf("constraint should not violate")
	}
	// flag with argument should skip next token
	if violatesPipPin([]string{"pip", "install", "--index-url", "u", "pkg"}) != true {
		t.Fatalf("expected violation with flag")
	}
	if violatesPipPin([]string{"python", "-m", "pip", "install", "pkg"}) != true {
		t.Fatalf("expected violation for python -m pip")
	}
}
