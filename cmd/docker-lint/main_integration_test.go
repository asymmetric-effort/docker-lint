// file: cmd/docker-lint/main_integration_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/rules"
	"github.com/asymmetric-effort/docker-lint/internal/version"
)

func TestIntegrationRunDetectsLatest(t *testing.T) {
	df := testDataPath("Dockerfile.bad")
	var out bytes.Buffer
	if err := run([]string{df}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	ids := map[string]struct{}{}
	for _, f := range findings {
		ids[f.RuleID] = struct{}{}
	}
	if _, ok := ids[rules.NewNoLatestTag().ID()]; !ok {
		t.Fatalf("missing DL3007 finding")
	}
	if _, ok := ids[rules.NewRequireOSVersionTag().ID()]; !ok {
		t.Fatalf("missing DL3043 finding")
	}
}

func TestIntegrationRunClean(t *testing.T) {
	df := testDataPath("Dockerfile.good")
	var out bytes.Buffer
	if err := run([]string{df}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRunNoArgs verifies that run returns a usage error when invoked with no arguments.
func TestIntegrationRunNoArgs(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{}, &out, io.Discard, false)
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

// TestIntegrationRunHelpShort verifies that run prints usage when the -h flag is provided.
func TestIntegrationRunHelpShort(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-h"}, &out, io.Discard, false); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "usage: docker-lint") {
		t.Fatalf("expected usage message, got %q", out.String())
	}
}

// TestIntegrationRunHelpLong verifies that run prints usage when the --help flag is provided.
func TestIntegrationRunHelpLong(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"--help"}, &out, io.Discard, false); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "usage: docker-lint") {
		t.Fatalf("expected usage message, got %q", out.String())
	}
}

// TestIntegrationRunMissingFile verifies that run returns an error when the Dockerfile is missing.
func TestIntegrationRunMissingFile(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"does-not-exist"}, &out, io.Discard, false)
	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected file not found error, got %v", err)
	}
}

// TestIntegrationRunInvalidDockerfile verifies that run returns an error for an invalid Dockerfile.
func TestIntegrationRunInvalidDockerfile(t *testing.T) {
	df := testDataPath("Dockerfile.invalid")
	var out bytes.Buffer
	if err := run([]string{df}, &out, io.Discard, false); err == nil {
		t.Fatalf("expected parse error, got nil")
	}
}

// TestIntegrationRunGlob verifies that run processes files matching single-star patterns.
func TestIntegrationRunGlob(t *testing.T) {
	tmp := t.TempDir()
	goodSrc := testDataPath("Dockerfile.good")
	badSrc := testDataPath("Dockerfile.bad")

	goodDst := filepath.Join(tmp, "Dockerfile.good")
	badDst := filepath.Join(tmp, "Dockerfile.bad")

	g, err := os.ReadFile(goodSrc)
	if err != nil {
		t.Fatalf("read good: %v", err)
	}
	if err := os.WriteFile(goodDst, g, 0o644); err != nil {
		t.Fatalf("write good: %v", err)
	}
	b, err := os.ReadFile(badSrc)
	if err != nil {
		t.Fatalf("read bad: %v", err)
	}
	if err := os.WriteFile(badDst, b, 0o644); err != nil {
		t.Fatalf("write bad: %v", err)
	}

	pattern := filepath.Join(tmp, "Dockerfile.*")
	var out bytes.Buffer
	if err := run([]string{pattern}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

// TestIntegrationRunDoubleStar verifies that run processes files matching recursive patterns.
func TestIntegrationRunDoubleStar(t *testing.T) {
	tmp := t.TempDir()
	nested := filepath.Join(tmp, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	badSrc := testDataPath("Dockerfile.bad")
	rootDst := filepath.Join(tmp, "Dockerfile.bad")
	nestedDst := filepath.Join(nested, "Dockerfile.bad")

	b, err := os.ReadFile(badSrc)
	if err != nil {
		t.Fatalf("read bad: %v", err)
	}
	if err := os.WriteFile(rootDst, b, 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}
	if err := os.WriteFile(nestedDst, b, 0o644); err != nil {
		t.Fatalf("write nested: %v", err)
	}

	pattern := filepath.Join(tmp, "**", "Dockerfile.bad")
	var out bytes.Buffer
	if err := run([]string{pattern}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 4 {
		t.Fatalf("expected 4 findings, got %d", len(findings))
	}
}

// TestIntegrationRunDefaultConfigIgnored verifies that .docker-lint.yaml ignored rules are applied.
func TestIntegrationRunDefaultConfigIgnored(t *testing.T) {
	tmp := t.TempDir()
	// write Dockerfile
	src := testDataPath("Dockerfile.bad")
	df := filepath.Join(tmp, "Dockerfile.bad")
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read dockerfile: %v", err)
	}
	if err := os.WriteFile(df, b, 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}
	// write config
	cfg := []byte("ignored:\n  - DL3007\n")
	if err := os.WriteFile(filepath.Join(tmp, ".docker-lint.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Chdir(tmp)
	var out bytes.Buffer
	if err := run([]string{"Dockerfile.bad"}, &out, io.Discard, false); err != nil {
		t.Fatalf("run: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != rules.NewRequireOSVersionTag().ID() {
		t.Fatalf("unexpected rule: %s", findings[0].RuleID)
	}
}

// TestIntegrationRunConfigFlagShort verifies that -c config flag applies ignored rules.
func TestIntegrationRunConfigFlagShort(t *testing.T) {
	tmp := t.TempDir()
	src := testDataPath("Dockerfile.bad")
	df := filepath.Join(tmp, "Dockerfile.bad")
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read dockerfile: %v", err)
	}
	if err := os.WriteFile(df, b, 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}
	cfgPath := filepath.Join(tmp, "cfg.yaml")
	cfg := []byte("ignored:\n  - DL3007\n")
	if err := os.WriteFile(cfgPath, cfg, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	var out bytes.Buffer
	if err := run([]string{"-c", cfgPath, df}, &out, io.Discard, false); err != nil {
		t.Fatalf("run: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != rules.NewRequireOSVersionTag().ID() {
		t.Fatalf("unexpected rule: %s", findings[0].RuleID)
	}
}

// TestIntegrationRunConfigFlagLong verifies that --config flag applies ignored rules.
func TestIntegrationRunConfigFlagLong(t *testing.T) {
	tmp := t.TempDir()
	src := testDataPath("Dockerfile.bad")
	df := filepath.Join(tmp, "Dockerfile.bad")
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read dockerfile: %v", err)
	}
	if err := os.WriteFile(df, b, 0o644); err != nil {
		t.Fatalf("write dockerfile: %v", err)
	}
	cfgPath := filepath.Join(tmp, "cfg.yaml")
	cfg := []byte("ignored:\n  - DL3007\n")
	if err := os.WriteFile(cfgPath, cfg, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	var out bytes.Buffer
	if err := run([]string{"--config", cfgPath, df}, &out, io.Discard, false); err != nil {
		t.Fatalf("run: %v", err)
	}
	var findings []engine.Finding
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != rules.NewRequireOSVersionTag().ID() {
		t.Fatalf("unexpected rule: %s", findings[0].RuleID)
	}
}

// TestIntegrationRunVersion ensures the version command prints the current version.
func TestIntegrationRunVersion(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"version"}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != version.Current {
		t.Fatalf("expected %q, got %q", version.Current, got)
	}
}

// TestIntegrationRunVersionFlagShort ensures -version prints the current version.
func TestIntegrationRunVersionFlagShort(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-version"}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != version.Current {
		t.Fatalf("expected %q, got %q", version.Current, got)
	}
}

// TestIntegrationRunVersionFlagLong ensures --version prints the current version.
func TestIntegrationRunVersionFlagLong(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"--version"}, &out, io.Discard, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != version.Current {
		t.Fatalf("expected %q, got %q", version.Current, got)
	}
}
