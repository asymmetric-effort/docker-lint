// file: internal/rules/DL1001_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationNoInlineIgnoreID validates rule identity.
func TestIntegrationNoInlineIgnoreID(t *testing.T) {
	if NewNoInlineIgnore().ID() != "DL1001" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationNoInlineIgnoreViolation reports inline ignore pragmas.
func TestIntegrationNoInlineIgnoreViolation(t *testing.T) {
	src := "# hadolint ignore=DL3007\nFROM alpine\nRUN echo hello # hadolint ignore=DL4000\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewNoInlineIgnore()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

// TestIntegrationNoInlineIgnoreClean ensures files without ignores pass.
func TestIntegrationNoInlineIgnoreClean(t *testing.T) {
	src := "FROM alpine\nRUN echo hello\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewNoInlineIgnore()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}
