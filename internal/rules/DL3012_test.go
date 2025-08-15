// file: internal/rules/DL3012_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationSingleHealthcheckID validates rule identity.
func TestIntegrationSingleHealthcheckID(t *testing.T) {
	if NewSingleHealthcheck().ID() != "DL3012" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationSingleHealthcheckViolation detects multiple healthchecks in a stage.
func TestIntegrationSingleHealthcheckViolation(t *testing.T) {
	src := "FROM alpine\nHEALTHCHECK CMD true\nHEALTHCHECK CMD false\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleHealthcheck()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 3 {
		t.Fatalf("expected one finding on line 3, got %#v", findings)
	}
}

// TestIntegrationSingleHealthcheckClean ensures single healthcheck passes.
func TestIntegrationSingleHealthcheckClean(t *testing.T) {
	src := "FROM alpine\nHEALTHCHECK CMD true\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleHealthcheck()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationSingleHealthcheckMultiStage allows one healthcheck per stage.
func TestIntegrationSingleHealthcheckMultiStage(t *testing.T) {
	src := "FROM alpine AS base\nHEALTHCHECK CMD true\nFROM scratch\nHEALTHCHECK CMD false\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleHealthcheck()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationSingleHealthcheckNilDocument ensures nil input is handled.
func TestIntegrationSingleHealthcheckNilDocument(t *testing.T) {
	r := NewSingleHealthcheck()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
