// file: internal/rules/DL4004_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationSingleEntrypointID validates rule identity.
func TestIntegrationSingleEntrypointID(t *testing.T) {
	if NewSingleEntrypoint().ID() != "DL4004" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationSingleEntrypointViolation detects multiple ENTRYPOINT instructions in a stage.
func TestIntegrationSingleEntrypointViolation(t *testing.T) {
	src := "FROM alpine\nENTRYPOINT [\"/bin/ls\"]\nENTRYPOINT [\"/bin/sh\"]\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleEntrypoint()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 3 {
		t.Fatalf("expected one finding on line 3, got %#v", findings)
	}
}

// TestIntegrationSingleEntrypointSeparateStages allows ENTRYPOINT per stage.
func TestIntegrationSingleEntrypointSeparateStages(t *testing.T) {
	src := "FROM alpine\nENTRYPOINT [\"/bin/ls\"]\nFROM ubuntu\nENTRYPOINT [\"/bin/bash\"]\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleEntrypoint()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationSingleEntrypointNilDocument ensures graceful handling of nil input.
func TestIntegrationSingleEntrypointNilDocument(t *testing.T) {
	r := NewSingleEntrypoint()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
