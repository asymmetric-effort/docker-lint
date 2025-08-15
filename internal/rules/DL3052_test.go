// file: internal/rules/DL3052_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationLabelURLValidID validates rule identity.
func TestIntegrationLabelURLValidID(t *testing.T) {
	if NewLabelURLValid(nil).ID() != "DL3052" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationLabelURLValidViolation detects malformed URLs.
func TestIntegrationLabelURLValidViolation(t *testing.T) {
	src := "FROM scratch\nLABEL homepage=not-a-url\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"homepage": LabelTypeURL}
	r := NewLabelURLValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationLabelURLValidClean ensures valid URLs pass.
func TestIntegrationLabelURLValidClean(t *testing.T) {
	src := "FROM scratch\nLABEL homepage=http://example.com\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"homepage": LabelTypeURL}
	r := NewLabelURLValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationLabelURLValidNilDocument ensures nil documents are handled gracefully.
func TestIntegrationLabelURLValidNilDocument(t *testing.T) {
	r := NewLabelURLValid(nil)
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
