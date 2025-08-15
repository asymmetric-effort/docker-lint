// file: internal/rules/DL4003_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationSingleCmdID validates rule identity.
func TestIntegrationSingleCmdID(t *testing.T) {
	if NewSingleCmd().ID() != "DL4003" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationSingleCmdViolation detects multiple CMD instructions in a stage.
func TestIntegrationSingleCmdViolation(t *testing.T) {
	src := "FROM alpine\nCMD echo one\nCMD echo two\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleCmd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 3 {
		t.Fatalf("expected one finding on line 3, got %#v", findings)
	}
}

// TestIntegrationSingleCmdClean ensures a single CMD passes.
func TestIntegrationSingleCmdClean(t *testing.T) {
	src := "FROM alpine\nCMD echo one\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleCmd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationSingleCmdMultiStage validates each stage may have one CMD.
func TestIntegrationSingleCmdMultiStage(t *testing.T) {
	src := "FROM alpine\nCMD echo one\nFROM alpine\nCMD echo two\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleCmd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationSingleCmdMultiStageViolation flags extra CMD in later stage.
func TestIntegrationSingleCmdMultiStageViolation(t *testing.T) {
	src := "FROM alpine\nCMD echo one\nFROM alpine\nCMD echo two\nCMD echo three\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleCmd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 5 {
		t.Fatalf("expected one finding on line 5, got %#v", findings)
	}
}

// TestIntegrationSingleCmdCaseInsensitive ensures lowercase cmd is detected.
func TestIntegrationSingleCmdCaseInsensitive(t *testing.T) {
	src := "FROM alpine\ncmd echo one\ncmd echo two\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewSingleCmd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 3 {
		t.Fatalf("expected one finding on line 3, got %#v", findings)
	}
}

// TestIntegrationSingleCmdNilDocument ensures graceful handling of nil input.
func TestIntegrationSingleCmdNilDocument(t *testing.T) {
	r := NewSingleCmd()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
