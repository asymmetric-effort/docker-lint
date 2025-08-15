// file: internal/rules/DL3022_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationCopyFromPreviousStageID validates rule identity.
func TestIntegrationCopyFromPreviousStageID(t *testing.T) {
	if NewCopyFromPreviousStage().ID() != "DL3022" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationCopyFromPreviousStageViolation detects missing aliases.
func TestIntegrationCopyFromPreviousStageViolation(t *testing.T) {
	src := "FROM alpine\nCOPY --from=bogus /src /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromPreviousStage()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationCopyFromPreviousStageNumeric validates numeric stage references.
func TestIntegrationCopyFromPreviousStageNumeric(t *testing.T) {
	src := "FROM alpine\nCOPY --from=1 /src /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromPreviousStage()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationCopyFromPreviousStageClean ensures compliant Dockerfiles pass.
func TestIntegrationCopyFromPreviousStageClean(t *testing.T) {
	src := "FROM alpine AS build\nFROM alpine\nCOPY --from=build /src /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromPreviousStage()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyFromPreviousStageExternalImage ignores external references.
func TestIntegrationCopyFromPreviousStageExternalImage(t *testing.T) {
	src := "FROM alpine\nCOPY --from=alpine:latest /src /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromPreviousStage()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyFromPreviousStageNil ensures graceful handling of nil input.
func TestIntegrationCopyFromPreviousStageNil(t *testing.T) {
	r := NewCopyFromPreviousStage()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
