// file: cmd/docker-lint/main_integration_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/rules"
)

func TestIntegrationRunDetectsLatest(t *testing.T) {
	df := testDataPath("Dockerfile.bad")
	var out bytes.Buffer
	if err := run([]string{df}, &out); err != nil {
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
	if err := run([]string{df}, &out); err != nil {
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
	err := run([]string{}, &out)
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

// TestIntegrationRunHelpShort verifies that run prints usage when the -h flag is provided.
func TestIntegrationRunHelpShort(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-h"}, &out); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "usage: docker-lint") {
		t.Fatalf("expected usage message, got %q", out.String())
	}
}

// TestIntegrationRunHelpLong verifies that run prints usage when the --help flag is provided.
func TestIntegrationRunHelpLong(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"--help"}, &out); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "usage: docker-lint") {
		t.Fatalf("expected usage message, got %q", out.String())
	}
}

// TestIntegrationRunMissingFile verifies that run returns an error when the Dockerfile is missing.
func TestIntegrationRunMissingFile(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"does-not-exist"}, &out)
	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected file not found error, got %v", err)
	}
}

// TestIntegrationRunInvalidDockerfile verifies that run returns an error for an invalid Dockerfile.
func TestIntegrationRunInvalidDockerfile(t *testing.T) {
	df := testDataPath("Dockerfile.invalid")
	var out bytes.Buffer
	if err := run([]string{df}, &out); err == nil {
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
	if err := run([]string{pattern}, &out); err != nil {
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
	if err := run([]string{pattern}, &out); err != nil {
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
