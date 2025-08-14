// file: cmd/docker-lint/main_integration_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"bytes"
	"encoding/json"
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
