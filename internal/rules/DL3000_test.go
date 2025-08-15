// file: internal/rules/DL3000_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationAbsoluteWorkdirID validates rule identity.
func TestIntegrationAbsoluteWorkdirID(t *testing.T) {
	if NewAbsoluteWorkdir().ID() != "DL3000" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationAbsoluteWorkdirViolation detects relative WORKDIR usage.
func TestIntegrationAbsoluteWorkdirViolation(t *testing.T) {
	src := "FROM alpine\nWORKDIR app\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAbsoluteWorkdir()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationAbsoluteWorkdirDotRelative detects relative paths with leading dot.
func TestIntegrationAbsoluteWorkdirDotRelative(t *testing.T) {
	src := "FROM alpine\nWORKDIR ./app\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAbsoluteWorkdir()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationAbsoluteWorkdirClean ensures absolute paths pass.
func TestIntegrationAbsoluteWorkdirClean(t *testing.T) {
	src := "FROM alpine\nWORKDIR /app\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAbsoluteWorkdir()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationAbsoluteWorkdirVariable allows variable-based paths.
func TestIntegrationAbsoluteWorkdirVariable(t *testing.T) {
	src := "FROM alpine\nWORKDIR $APPDIR\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAbsoluteWorkdir()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationAbsoluteWorkdirWindows accepts Windows-style paths.
func TestIntegrationAbsoluteWorkdirWindows(t *testing.T) {
	src := "FROM alpine\nWORKDIR C:\\\build\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAbsoluteWorkdir()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationAbsoluteWorkdirNilDocument ensures nil input is handled.
func TestIntegrationAbsoluteWorkdirNilDocument(t *testing.T) {
	r := NewAbsoluteWorkdir()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
