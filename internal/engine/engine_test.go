// file: internal/engine/engine_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package engine_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	engine "github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"github.com/asymmetric-effort/docker-lint/internal/rules"
)

type stubRule struct {
	id       string
	findings []engine.Finding
	err      error
}

func (s stubRule) ID() string { return s.id }

func (s stubRule) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	return s.findings, s.err
}

// TestIntegrationRegistryRun verifies successful rule execution.
func TestIntegrationRegistryRun(t *testing.T) {
	r := engine.NewRegistry()
	r.Register(stubRule{id: "A", findings: []engine.Finding{{RuleID: "A"}}})
	out, err := r.Run(context.Background(), &ir.Document{})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if len(out) != 1 || out[0].RuleID != "A" {
		t.Fatalf("unexpected findings: %#v", out)
	}
}

// TestIntegrationRegistryRunError ensures errors propagate from rules.
func TestIntegrationRegistryRunError(t *testing.T) {
	r := engine.NewRegistry()
	r.Register(stubRule{id: "A", err: errors.New("bad")})
	if _, err := r.Run(context.Background(), &ir.Document{}); err == nil {
		t.Fatalf("expected error")
	}
}

// TestIntegrationInlineIgnorePrevLine ensures preceding ignore pragmas skip rules.
func TestIntegrationInlineIgnorePrevLine(t *testing.T) {
	src := "# hadolint ignore=DL3007\nFROM alpine\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := engine.NewRegistry()
	r.Register(rules.NewNoLatestTag())
	findings, err := r.Run(context.Background(), doc)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationInlineIgnoreSameLine ensures trailing ignore pragmas skip rules.
func TestIntegrationInlineIgnoreSameLine(t *testing.T) {
	src := "FROM alpine # hadolint ignore=DL3007\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := engine.NewRegistry()
	r.Register(rules.NewNoLatestTag())
	findings, err := r.Run(context.Background(), doc)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}
