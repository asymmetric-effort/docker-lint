// file: internal/rules/DL3009_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationAptListsCleanupID validates rule identity.
func TestIntegrationAptListsCleanupID(t *testing.T) {
	if NewAptListsCleanup().ID() != "DL3009" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationAptListsCleanupViolation detects missing cleanup after apt install.
func TestIntegrationAptListsCleanupViolation(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAptListsCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationAptListsCleanupClean ensures cleanup after install passes.
func TestIntegrationAptListsCleanupClean(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAptListsCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationAptListsCleanupCleanOnly verifies apt-get clean is insufficient.
func TestIntegrationAptListsCleanupCleanOnly(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl && apt-get clean\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAptListsCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationAptListsCleanupNil ensures nil documents are handled.
func TestIntegrationAptListsCleanupNil(t *testing.T) {
	r := NewAptListsCleanup()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
