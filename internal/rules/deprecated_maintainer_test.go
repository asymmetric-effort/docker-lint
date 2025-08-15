// file: internal/rules/deprecated_maintainer_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationDeprecatedMaintainerID validates rule identity.
func TestIntegrationDeprecatedMaintainerID(t *testing.T) {
	if NewDeprecatedMaintainer().ID() != "DL4000" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationDeprecatedMaintainerViolation detects deprecated MAINTAINER usage.
func TestIntegrationDeprecatedMaintainerViolation(t *testing.T) {
	src := "FROM alpine\nMAINTAINER Somebody\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDeprecatedMaintainer()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationDeprecatedMaintainerClean ensures compliant Dockerfiles pass.
func TestIntegrationDeprecatedMaintainerClean(t *testing.T) {
	src := "FROM alpine\nLABEL maintainer=\"Somebody\"\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDeprecatedMaintainer()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationDeprecatedMaintainerCaseInsensitive catches lowercase usage.
func TestIntegrationDeprecatedMaintainerCaseInsensitive(t *testing.T) {
	src := "FROM alpine\nmaintainer Somebody\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDeprecatedMaintainer()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationDeprecatedMaintainerNilDocument ensures graceful handling of nil input.
func TestIntegrationDeprecatedMaintainerNilDocument(t *testing.T) {
	r := NewDeprecatedMaintainer()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
