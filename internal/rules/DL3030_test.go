// file: internal/rules/DL3030_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationRequireYumYesID validates rule identity.
func TestIntegrationRequireYumYesID(t *testing.T) {
	if NewRequireYumYes().ID() != "DL3030" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationRequireYumYesViolation detects missing -y.
func TestIntegrationRequireYumYesViolation(t *testing.T) {
	src := "FROM alpine\nRUN yum install wget\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewRequireYumYes()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationRequireYumYesClean ensures compliant commands pass.
func TestIntegrationRequireYumYesClean(t *testing.T) {
	src := "FROM alpine\nRUN yum install -y wget\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewRequireYumYes()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireYumYesNil ensures graceful nil handling.
func TestIntegrationRequireYumYesNil(t *testing.T) {
	r := NewRequireYumYes()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
