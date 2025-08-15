// file: internal/config/config_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoad verifies that Load parses hadolint-style configuration files.
func TestLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cfg.yaml")
	src := []byte("" +
		"ignored:\n" +
		"  - DL3007\n" +
		"override:\n" +
		"  warning:\n" +
		"    - SC1099\n" +
		"failure-threshold: warning\n" +
		"trustedRegistries:\n" +
		"  - ghcr.io\n" +
		"strict-labels: true\n" +
		"label-schema:\n" +
		"  author: text\n")
	if err := os.WriteFile(path, src, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(cfg.Ignored) != 1 || cfg.Ignored[0] != "DL3007" {
		t.Fatalf("unexpected ignored: %v", cfg.Ignored)
	}
	if cfg.FailureThreshold != "warning" {
		t.Fatalf("unexpected failure threshold: %s", cfg.FailureThreshold)
	}
	if !cfg.StrictLabels {
		t.Fatalf("expected strict labels enabled")
	}
	if v := cfg.LabelSchema["author"]; v != "text" {
		t.Fatalf("unexpected label schema: %v", cfg.LabelSchema)
	}
	if r := cfg.Override["warning"]; len(r) != 1 || r[0] != "SC1099" {
		t.Fatalf("unexpected override: %v", cfg.Override)
	}
	if len(cfg.TrustedRegistries) != 1 || cfg.TrustedRegistries[0] != "ghcr.io" {
		t.Fatalf("unexpected registries: %v", cfg.TrustedRegistries)
	}
}

// TestIsIgnored verifies that IsIgnored returns true for configured rules.
func TestIsIgnored(t *testing.T) {
	cfg := Config{Ignored: []string{"DL3007"}}
	if !cfg.IsIgnored("DL3007") {
		t.Fatalf("expected rule ignored")
	}
	if cfg.IsIgnored("DL3043") {
		t.Fatalf("unexpected ignore")
	}
}
