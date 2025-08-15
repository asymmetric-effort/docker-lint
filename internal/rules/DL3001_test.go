// file: internal/rules/DL3001_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationNoIrrelevantCommandsID validates rule identity.
func TestIntegrationNoIrrelevantCommandsID(t *testing.T) {
	if NewNoIrrelevantCommands().ID() != "DL3001" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationNoIrrelevantCommandsViolation detects banned commands.
func TestIntegrationNoIrrelevantCommandsViolation(t *testing.T) {
	src := "FROM alpine\nRUN ssh localhost\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewNoIrrelevantCommands()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationNoIrrelevantCommandsExecForm detects banned commands in exec form.
func TestIntegrationNoIrrelevantCommandsExecForm(t *testing.T) {
	src := "FROM alpine\nRUN [\"ssh\",\"localhost\"]\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewNoIrrelevantCommands()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationNoIrrelevantCommandsConnector handles multiple commands.
func TestIntegrationNoIrrelevantCommandsConnector(t *testing.T) {
	src := "FROM alpine\nRUN echo hi && ssh localhost\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewNoIrrelevantCommands()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationNoIrrelevantCommandsClean ensures compliant Dockerfiles pass.
func TestIntegrationNoIrrelevantCommandsClean(t *testing.T) {
	src := "FROM alpine\nRUN echo hi\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewNoIrrelevantCommands()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationNoIrrelevantCommandsNilDocument ensures graceful handling of nil input.
func TestIntegrationNoIrrelevantCommandsNilDocument(t *testing.T) {
	r := NewNoIrrelevantCommands()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
