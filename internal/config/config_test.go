// file: internal/config/config_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoad verifies that Load reads exclusions from a YAML file.
func TestLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cfg.yaml")
	data := []byte("exclusions:\n  - DL3007\n  - DL3043\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(cfg.Exclusions) != 2 {
		t.Fatalf("expected 2 exclusions, got %d", len(cfg.Exclusions))
	}
	if cfg.Exclusions[0] != "DL3007" || cfg.Exclusions[1] != "DL3043" {
		t.Fatalf("unexpected exclusions: %v", cfg.Exclusions)
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
