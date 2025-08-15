// file: internal/rules/DL3037_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationPinZypperVersionsID validates rule identity.
func TestIntegrationPinZypperVersionsID(t *testing.T) {
	if NewPinZypperVersions().ID() != "DL3037" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationPinZypperVersionsViolation detects unpinned packages.
func TestIntegrationPinZypperVersionsViolation(t *testing.T) {
	src := "FROM alpine\nRUN zypper install pkg\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinZypperVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationPinZypperVersionsClean ensures versioned packages pass.
func TestIntegrationPinZypperVersionsClean(t *testing.T) {
	src := "FROM alpine\nRUN zypper install pkg=1.0\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinZypperVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPinZypperVersionsNil ensures graceful nil handling.
func TestIntegrationPinZypperVersionsNil(t *testing.T) {
	r := NewPinZypperVersions()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
