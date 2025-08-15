// file: internal/rules/DL3002_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationLastUserNotRootID validates rule identity.
func TestIntegrationLastUserNotRootID(t *testing.T) {
	if NewLastUserNotRoot().ID() != "DL3002" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationLastUserNotRootViolation detects stages ending with root user.
func TestIntegrationLastUserNotRootViolation(t *testing.T) {
	src := "FROM alpine\nUSER root\nRUN echo hi\nFROM busybox\nUSER 0:0\nFROM scratch\nUSER app\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewLastUserNotRoot()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 || findings[0].Line != 2 || findings[1].Line != 5 {
		t.Fatalf("expected findings on lines 2 and 5, got %#v", findings)
	}
}

// TestIntegrationLastUserNotRootGroup detects root with group specification.
func TestIntegrationLastUserNotRootGroup(t *testing.T) {
	src := "FROM alpine\nUSER root:root\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewLastUserNotRoot()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationLastUserNotRootClean ensures compliant Dockerfiles pass.
func TestIntegrationLastUserNotRootClean(t *testing.T) {
	src := "FROM alpine\nUSER root\nUSER app\nRUN echo hi\nFROM busybox\nRUN echo hi\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewLastUserNotRoot()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationLastUserNotRootNilDocument ensures graceful handling of nil input.
func TestIntegrationLastUserNotRootNilDocument(t *testing.T) {
	r := NewLastUserNotRoot()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
