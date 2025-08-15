// file: internal/config/config_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadAndIsRuleExcluded verifies configuration loading and exclusion checks.
func TestLoadAndIsRuleExcluded(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".docker-lint.yaml")
	src := "exclude:\n  Dockerfile.bad:\n    - DL3007\n"
	if err := os.WriteFile(cfgPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !cfg.IsRuleExcluded(filepath.Join(tmp, "Dockerfile.bad"), "DL3007") {
		t.Fatalf("expected rule excluded")
	}
	if cfg.IsRuleExcluded(filepath.Join(tmp, "Dockerfile.bad"), "DL3043") {
		t.Fatalf("unexpected exclusion")
	}
}
