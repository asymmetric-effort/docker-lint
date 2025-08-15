// file: internal/rules/DL3058_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

func TestLabelEmailValidID(t *testing.T) {
	if NewLabelEmailValid(nil).ID() != "DL3058" {
		t.Fatalf("unexpected id")
	}
}

func TestLabelEmailValidViolation(t *testing.T) {
	src := "FROM scratch\nLABEL maintainer=not-an-email\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"maintainer": LabelTypeEmail}
	r := NewLabelEmailValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

func TestLabelEmailValidClean(t *testing.T) {
	src := "FROM scratch\nLABEL maintainer=test@example.com\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"maintainer": LabelTypeEmail}
	r := NewLabelEmailValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestLabelEmailValidNilDocument(t *testing.T) {
	r := NewLabelEmailValid(nil)
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
