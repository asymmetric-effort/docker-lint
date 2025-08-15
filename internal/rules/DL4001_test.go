// file: internal/rules/DL4001_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationExclusiveCurlWgetID validates rule identity.
func TestIntegrationExclusiveCurlWgetID(t *testing.T) {
	if NewExclusiveCurlWget().ID() != "DL4001" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationExclusiveCurlWgetViolation detects mixing curl and wget across RUNs.
func TestIntegrationExclusiveCurlWgetViolation(t *testing.T) {
	src := "FROM alpine\nRUN curl -L http://example.com\nRUN wget http://example.com\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewExclusiveCurlWget()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 3 {
		t.Fatalf("expected one finding on line 3, got %#v", findings)
	}
}

// TestIntegrationExclusiveCurlWgetSameRun flags using curl and wget in a single RUN.
func TestIntegrationExclusiveCurlWgetSameRun(t *testing.T) {
	src := "FROM alpine\nRUN wget http://a && curl http://b\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewExclusiveCurlWget()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationExclusiveCurlWgetCleanSingleTool allows repeated use of one tool.
func TestIntegrationExclusiveCurlWgetCleanSingleTool(t *testing.T) {
	src := "FROM alpine\nRUN curl -L http://a\nRUN curl -L http://b\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewExclusiveCurlWget()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationExclusiveCurlWgetCleanSeparateStages permits different tools per stage.
func TestIntegrationExclusiveCurlWgetCleanSeparateStages(t *testing.T) {
	src := "FROM alpine\nRUN curl -L http://a\nFROM alpine\nRUN wget http://b\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewExclusiveCurlWget()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationExclusiveCurlWgetNilDocument ensures graceful handling of nil input.
func TestIntegrationExclusiveCurlWgetNilDocument(t *testing.T) {
	r := NewExclusiveCurlWget()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
