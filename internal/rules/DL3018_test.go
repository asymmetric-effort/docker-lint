// file: internal/rules/DL3018_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationApkPinID validates rule identity.
func TestIntegrationApkPinID(t *testing.T) {
	if NewApkPin().ID() != "DL3018" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationApkPinViolation detects unpinned apk adds.
func TestIntegrationApkPinViolation(t *testing.T) {
	r := NewApkPin()
	src := "FROM alpine\nRUN apk add curl\n"
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

// TestIntegrationApkPinClean ensures compliant apk adds pass.
func TestIntegrationApkPinClean(t *testing.T) {
	r := NewApkPin()
	src := "FROM alpine\nRUN apk add curl=8.0.1 bash=5.1.0\n"
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

// TestIntegrationApkPinApkFile ensures .apk packages are treated as pinned.
func TestIntegrationApkPinApkFile(t *testing.T) {
	r := NewApkPin()
	src := "FROM alpine\nRUN apk add /tmp/pkg.apk\n"
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

// TestIntegrationApkPinNilDocument ensures graceful handling of nil input.
func TestIntegrationApkPinNilDocument(t *testing.T) {
	r := NewApkPin()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
