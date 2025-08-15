// file: internal/rules/DL3008_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationAptPinID validates rule identity.
func TestIntegrationAptPinID(t *testing.T) {
	if NewAptPin().ID() != "DL3008" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationAptPinViolation detects unpinned apt installs.
func TestIntegrationAptPinViolation(t *testing.T) {
	r := NewAptPin()
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl ca-certificates\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

// TestIntegrationAptPinMissingVersion detects packages with empty version.
func TestIntegrationAptPinMissingVersion(t *testing.T) {
	r := NewAptPin()
	src := "FROM ubuntu\nRUN apt-get install curl= ca-certificates=1.2\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

// TestIntegrationAptPinClean ensures compliant installs pass.
func TestIntegrationAptPinClean(t *testing.T) {
	r := NewAptPin()
	src := "FROM ubuntu\nRUN apt-get update && apt-get install -y curl=7.81.0-1 ca-certificates=20240203\nRUN apt install python3=3.10.*\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationAptPinEmptyDocument ensures empty documents pass.
func TestIntegrationAptPinEmptyDocument(t *testing.T) {
	r := NewAptPin()
	findings, err := r.Check(context.Background(), &ir.Document{})
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}
