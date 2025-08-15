// file: internal/rules/DL4006_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationPipefailBeforePipeID validates rule identity.
func TestIntegrationPipefailBeforePipeID(t *testing.T) {
	if NewPipefailBeforePipe().ID() != "DL4006" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationPipefailBeforePipeViolation warns on pipelines without pipefail.
func TestIntegrationPipefailBeforePipeViolation(t *testing.T) {
	src := "FROM alpine\nRUN echo hi | grep h\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPipefailBeforePipe()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationPipefailBeforePipeClean ensures pipefail option suppresses warning.
func TestIntegrationPipefailBeforePipeClean(t *testing.T) {
	src := "FROM alpine\nSHELL [\"/bin/bash\",\"-o\",\"pipefail\",\"-c\"]\nRUN echo hi | grep h\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPipefailBeforePipe()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPipefailBeforePipeNoPipe ignores commands without pipes.
func TestIntegrationPipefailBeforePipeNoPipe(t *testing.T) {
	src := "FROM alpine\nRUN echo hi && grep h file\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPipefailBeforePipe()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPipefailBeforePipeNonPosixShell ignores non-POSIX shells.
func TestIntegrationPipefailBeforePipeNonPosixShell(t *testing.T) {
	src := "FROM scratch\nSHELL [\"pwsh\",\"-c\"]\nRUN Get-Item a | Select-Object\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPipefailBeforePipe()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPipefailBeforePipeReset verifies state resets on new FROM.
func TestIntegrationPipefailBeforePipeReset(t *testing.T) {
	src := "FROM alpine\nSHELL [\"/bin/bash\",\"-o\",\"pipefail\",\"-c\"]\nRUN echo hi | grep h\nFROM alpine\nRUN echo hi | grep h\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewPipefailBeforePipe()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 5 {
		t.Fatalf("expected one finding on line 5, got %#v", findings)
	}
}

// TestIntegrationPipefailBeforePipeNil handles nil documents gracefully.
func TestIntegrationPipefailBeforePipeNil(t *testing.T) {
	r := NewPipefailBeforePipe()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
