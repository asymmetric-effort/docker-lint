// file: internal/engine/engine_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

type stubRule struct {
	id       string
	findings []Finding
	err      error
}

func (s stubRule) ID() string { return s.id }

func (s stubRule) Check(ctx context.Context, d *ir.Document) ([]Finding, error) {
	return s.findings, s.err
}

// TestIntegrationRegistryRun verifies successful rule execution.
func TestIntegrationRegistryRun(t *testing.T) {
	r := NewRegistry()
	r.Register(stubRule{id: "A", findings: []Finding{{RuleID: "A"}}})
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
	r := NewRegistry()
	r.Register(stubRule{id: "A", err: errors.New("bad")})
	if _, err := r.Run(context.Background(), &ir.Document{}); err == nil {
		t.Fatalf("expected error")
	}
}
