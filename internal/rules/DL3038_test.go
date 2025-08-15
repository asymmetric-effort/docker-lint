// file: internal/rules/DL3038_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationRequireDnfYesID validates rule identity.
func TestIntegrationRequireDnfYesID(t *testing.T) {
	if NewRequireDnfYes().ID() != "DL3038" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationRequireDnfYesViolation detects missing -y.
func TestIntegrationRequireDnfYesViolation(t *testing.T) {
	src := "FROM alpine\nRUN dnf install wget\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewRequireDnfYes()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationRequireDnfYesClean ensures compliant commands pass.
func TestIntegrationRequireDnfYesClean(t *testing.T) {
	src := "FROM alpine\nRUN dnf install -y wget\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewRequireDnfYes()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireDnfYesNil ensures graceful nil handling.
func TestIntegrationRequireDnfYesNil(t *testing.T) {
	r := NewRequireDnfYes()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
