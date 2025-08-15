// file: internal/rules/DL3025_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationJSONNotationCmdEntrypointID validates rule identity.
func TestIntegrationJSONNotationCmdEntrypointID(t *testing.T) {
	if NewJSONNotationCmdEntrypoint().ID() != "DL3025" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationJSONNotationCmdEntrypointViolation detects shell form usage.
func TestIntegrationJSONNotationCmdEntrypointViolation(t *testing.T) {
	src := "FROM alpine\nCMD echo hi\nENTRYPOINT echo hi\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewJSONNotationCmdEntrypoint()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected two findings, got %d", len(findings))
	}
}

// TestIntegrationJSONNotationCmdEntrypointClean ensures compliant Dockerfiles pass.
func TestIntegrationJSONNotationCmdEntrypointClean(t *testing.T) {
	src := "FROM alpine\nCMD [\"echo\",\"hi\"]\nENTRYPOINT [\"echo\",\"hi\"]\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewJSONNotationCmdEntrypoint()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationJSONNotationCmdEntrypointNil ensures graceful handling of nil input.
func TestIntegrationJSONNotationCmdEntrypointNil(t *testing.T) {
	r := NewJSONNotationCmdEntrypoint()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
