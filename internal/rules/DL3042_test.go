// file: internal/rules/DL3042_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationCombinePackageRunsID validates rule identity.
func TestIntegrationCombinePackageRunsID(t *testing.T) {
	if NewCombinePackageRuns().ID() != "DL3042" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationCombinePackageRunsViolation detects separated apt-get operations.
func TestIntegrationCombinePackageRunsViolation(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update\nRUN apt-get install -y curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCombinePackageRuns()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationCombinePackageRunsDifferentManagers ensures different managers are allowed.
func TestIntegrationCombinePackageRunsDifferentManagers(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl\nRUN apk add bash\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCombinePackageRuns()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCombinePackageRunsBreak ensures unrelated RUN breaks sequence.
func TestIntegrationCombinePackageRunsBreak(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update\nRUN echo ok\nRUN apt-get install -y curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCombinePackageRuns()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCombinePackageRunsClean ensures combined operations pass.
func TestIntegrationCombinePackageRunsClean(t *testing.T) {
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCombinePackageRuns()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCombinePackageRunsNil ensures nil documents are handled.
func TestIntegrationCombinePackageRunsNil(t *testing.T) {
	r := NewCombinePackageRuns()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
