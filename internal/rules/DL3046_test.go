// file: internal/rules/DL3046_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestApkNoUpgradeID verifies rule identity.
func TestApkNoUpgradeID(t *testing.T) {
	if NewApkNoUpgrade().ID() != "DL3046" {
		t.Fatalf("unexpected id")
	}
}

// TestApkNoUpgradeViolation flags apk upgrade usage.
func TestApkNoUpgradeViolation(t *testing.T) {
	r := NewApkNoUpgrade()
	src := "FROM alpine\nRUN apk --no-cache upgrade\n"
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

// TestApkNoUpgradeClean ensures other apk operations pass.
func TestApkNoUpgradeClean(t *testing.T) {
	r := NewApkNoUpgrade()
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
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestApkNoUpgradeNilDocument ensures nil or empty documents are handled.
func TestApkNoUpgradeNilDocument(t *testing.T) {
	r := NewApkNoUpgrade()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
