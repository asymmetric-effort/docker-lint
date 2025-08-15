// file: internal/rules/DL3024_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationUniqueStageNamesID validates rule identity.
func TestIntegrationUniqueStageNamesID(t *testing.T) {
	if NewUniqueStageNames().ID() != "DL3024" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationUniqueStageNamesViolation detects duplicate aliases.
func TestIntegrationUniqueStageNamesViolation(t *testing.T) {
	src := "FROM alpine AS base\nFROM ubuntu AS base\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewUniqueStageNames()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationUniqueStageNamesClean ensures compliant Dockerfiles pass.
func TestIntegrationUniqueStageNamesClean(t *testing.T) {
	src := "FROM alpine AS base\nFROM ubuntu AS build\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewUniqueStageNames()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUniqueStageNamesNil ensures graceful handling of nil input.
func TestIntegrationUniqueStageNamesNil(t *testing.T) {
	r := NewUniqueStageNames()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
