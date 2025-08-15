// file: internal/rules/DL3048_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationLabelKeyValidID validates rule identity.
func TestIntegrationLabelKeyValidID(t *testing.T) {
	if NewLabelKeyValid().ID() != "DL3048" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationLabelKeyValidViolation reports invalid label keys.
func TestIntegrationLabelKeyValidViolation(t *testing.T) {
	src := "FROM scratch\nLABEL com.docker.foo=bar invalid$key=x valid-label=yes\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewLabelKeyValid()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected two findings, got %d", len(findings))
	}
}

// TestIntegrationLabelKeyValidOnbuild detects issues in ONBUILD LABEL instructions.
func TestIntegrationLabelKeyValidOnbuild(t *testing.T) {
	src := "FROM scratch\nONBUILD LABEL some_key=value\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewLabelKeyValid()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationLabelKeyValidClean ensures compliant Dockerfiles pass.
func TestIntegrationLabelKeyValidClean(t *testing.T) {
	src := "FROM scratch\nLABEL org.example.meta.build=2025-08-14\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewLabelKeyValid()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationLabelKeyValidNilDocument ensures nil documents are handled gracefully.
func TestIntegrationLabelKeyValidNilDocument(t *testing.T) {
	r := NewLabelKeyValid()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
