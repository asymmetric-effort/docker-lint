// file: internal/rules/DL3033_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationPinYumVersionsID validates rule identity.
func TestIntegrationPinYumVersionsID(t *testing.T) {
	if NewPinYumVersions().ID() != "DL3033" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationPinYumVersionsViolation detects unpinned packages.
func TestIntegrationPinYumVersionsViolation(t *testing.T) {
	src := "FROM alpine\nRUN yum install pkg\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinYumVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationPinYumVersionsClean ensures versioned packages pass.
func TestIntegrationPinYumVersionsClean(t *testing.T) {
	src := "FROM alpine\nRUN yum install pkg-1.0\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinYumVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPinYumVersionsModule ensures modules require version.
func TestIntegrationPinYumVersionsModule(t *testing.T) {
	src := "FROM alpine\nRUN yum module install nodejs\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPinYumVersions()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationPinYumVersionsNil ensures graceful handling of nil input.
func TestIntegrationPinYumVersionsNil(t *testing.T) {
	r := NewPinYumVersions()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
