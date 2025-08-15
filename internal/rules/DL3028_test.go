// file: internal/rules/DL3028_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationPinGemVersionsID validates rule identity.
func TestIntegrationPinGemVersionsID(t *testing.T) {
	if NewPinGemVersions().ID() != "DL3028" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationPinGemVersionsViolation detects unpinned gem installs.
func TestIntegrationPinGemVersionsViolation(t *testing.T) {
	src := "FROM alpine\nRUN gem install rake\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinGemVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationPinGemVersionsClean ensures compliant Dockerfiles pass.
func TestIntegrationPinGemVersionsClean(t *testing.T) {
	src := "FROM alpine\nRUN gem install rake:1.0\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinGemVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPinGemVersionsOptionSkipped ensures version options skip rule.
func TestIntegrationPinGemVersionsOptionSkipped(t *testing.T) {
	src := "FROM alpine\nRUN gem install rake -v 1.2.3\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinGemVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPinGemVersionsNil ensures graceful handling of nil input.
func TestIntegrationPinGemVersionsNil(t *testing.T) {
	r := NewPinGemVersions()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
