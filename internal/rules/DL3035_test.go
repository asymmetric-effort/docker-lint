// file: internal/rules/DL3035_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationForbidZypperDistUpgradeID validates rule identity.
func TestIntegrationForbidZypperDistUpgradeID(t *testing.T) {
	if NewForbidZypperDistUpgrade().ID() != "DL3035" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationForbidZypperDistUpgradeViolation detects disallowed command.
func TestIntegrationForbidZypperDistUpgradeViolation(t *testing.T) {
	src := "FROM alpine\nRUN zypper dist-upgrade\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewForbidZypperDistUpgrade()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationForbidZypperDistUpgradeClean ensures other commands pass.
func TestIntegrationForbidZypperDistUpgradeClean(t *testing.T) {
	src := "FROM alpine\nRUN zypper install -y pkg\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewForbidZypperDistUpgrade()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationForbidZypperDistUpgradeNil ensures graceful nil handling.
func TestIntegrationForbidZypperDistUpgradeNil(t *testing.T) {
	r := NewForbidZypperDistUpgrade()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
