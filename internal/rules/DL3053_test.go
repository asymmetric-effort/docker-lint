// file: internal/rules/DL3053_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

func TestLabelTimeRFC3339ID(t *testing.T) {
	if NewLabelTimeRFC3339(nil).ID() != "DL3053" {
		t.Fatalf("unexpected id")
	}
}

func TestLabelTimeRFC3339Violation(t *testing.T) {
	src := "FROM scratch\nLABEL built=not-time\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"built": LabelTypeRFC3339}
	r := NewLabelTimeRFC3339(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

func TestLabelTimeRFC3339Clean(t *testing.T) {
	src := "FROM scratch\nLABEL built=2025-01-01T00:00:00Z\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"built": LabelTypeRFC3339}
	r := NewLabelTimeRFC3339(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestLabelTimeRFC3339NilDocument(t *testing.T) {
	r := NewLabelTimeRFC3339(nil)
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
