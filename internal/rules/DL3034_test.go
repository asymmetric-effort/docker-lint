// file: internal/rules/DL3034_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationRequireZypperYesID validates rule identity.
func TestIntegrationRequireZypperYesID(t *testing.T) {
	if NewRequireZypperYes().ID() != "DL3034" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationRequireZypperYesViolation detects missing flags.
func TestIntegrationRequireZypperYesViolation(t *testing.T) {
	src := "FROM alpine\nRUN zypper install pkg\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewRequireZypperYes()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationRequireZypperYesClean ensures non-interactive flag passes.
func TestIntegrationRequireZypperYesClean(t *testing.T) {
	src := "FROM alpine\nRUN zypper install -y pkg\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewRequireZypperYes()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireZypperYesNil ensures graceful nil handling.
func TestIntegrationRequireZypperYesNil(t *testing.T) {
	r := NewRequireZypperYes()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
