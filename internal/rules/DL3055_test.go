// file: internal/rules/DL3055_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

func TestLabelGitHashValidID(t *testing.T) {
	if NewLabelGitHashValid(nil).ID() != "DL3055" {
		t.Fatalf("unexpected id")
	}
}

func TestLabelGitHashValidViolation(t *testing.T) {
	src := "FROM scratch\nLABEL commit=xyz\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"commit": LabelTypeGitHash}
	r := NewLabelGitHashValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

func TestLabelGitHashValidClean(t *testing.T) {
	src := "FROM scratch\nLABEL commit=0123456789abcdef0123456789abcdef01234567\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"commit": LabelTypeGitHash}
	r := NewLabelGitHashValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestLabelGitHashValidNilDocument(t *testing.T) {
	r := NewLabelGitHashValid(nil)
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
