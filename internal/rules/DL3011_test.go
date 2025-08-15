// file: internal/rules/DL3011_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationValidPortRangeID validates rule identity.
func TestIntegrationValidPortRangeID(t *testing.T) {
	if NewValidPortRange().ID() != "DL3011" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationValidPortRangeViolation detects invalid port declarations.
func TestIntegrationValidPortRangeViolation(t *testing.T) {
	src := "FROM alpine\nEXPOSE 80 65536\nEXPOSE 8080-90000/tcp\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewValidPortRange()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 || findings[0].Line != 2 || findings[1].Line != 3 {
		t.Fatalf("expected findings on lines 2 and 3, got %#v", findings)
	}
}

// TestIntegrationValidPortRangeClean ensures compliant Dockerfiles pass.
func TestIntegrationValidPortRangeClean(t *testing.T) {
	src := "FROM alpine\nEXPOSE 80 443/tcp\nEXPOSE 1000-2000/udp\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewValidPortRange()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationValidPortRangeNilDocument ensures graceful handling of nil input.
func TestIntegrationValidPortRangeNilDocument(t *testing.T) {
	r := NewValidPortRange()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
