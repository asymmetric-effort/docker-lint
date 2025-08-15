// file: internal/rules/DL4005_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationUseShellForDefaultID validates rule identity.
func TestIntegrationUseShellForDefaultID(t *testing.T) {
	if NewUseShellForDefault().ID() != "DL4005" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationUseShellForDefaultViolation detects linking /bin/sh via RUN.
func TestIntegrationUseShellForDefaultViolation(t *testing.T) {
	src := "FROM alpine\nRUN ln -sf /bin/bash /bin/sh\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseShellForDefault()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationUseShellForDefaultClean ensures proper SHELL usage passes.
func TestIntegrationUseShellForDefaultClean(t *testing.T) {
	src := "FROM alpine\nSHELL [\"/bin/bash\", \"-c\"]\nRUN echo hi\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseShellForDefault()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUseShellForDefaultChained detects ln within command chains.
func TestIntegrationUseShellForDefaultChained(t *testing.T) {
	src := "FROM alpine\nRUN echo hi && ln -s /bin/bash /bin/sh\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseShellForDefault()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationUseShellForDefaultNilDocument ensures nil documents are handled.
func TestIntegrationUseShellForDefaultNilDocument(t *testing.T) {
	r := NewUseShellForDefault()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
