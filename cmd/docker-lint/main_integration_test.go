// file: cmd/docker-lint/main_integration_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
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
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != rules.NewNoLatestTag().ID() {
		t.Fatalf("unexpected rule id: %s", findings[0].RuleID)
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
