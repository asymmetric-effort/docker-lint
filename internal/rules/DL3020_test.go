// file: internal/rules/DL3020_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationUseCopyInsteadOfAddID validates rule identity.
func TestIntegrationUseCopyInsteadOfAddID(t *testing.T) {
	if NewUseCopyInsteadOfAdd().ID() != "DL3020" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationUseCopyInsteadOfAddViolation reports using ADD for local files.
func TestIntegrationUseCopyInsteadOfAddViolation(t *testing.T) {
	src := "FROM alpine\nADD file /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseCopyInsteadOfAdd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 2 {
		t.Fatalf("expected one finding on line 2, got %#v", findings)
	}
}

// TestIntegrationUseCopyInsteadOfAddRemote allows remote URLs.
func TestIntegrationUseCopyInsteadOfAddRemote(t *testing.T) {
	src := "FROM alpine\nADD https://example.com/file.tgz /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseCopyInsteadOfAdd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUseCopyInsteadOfAddArchive allows archive extraction.
func TestIntegrationUseCopyInsteadOfAddArchive(t *testing.T) {
	src := "FROM alpine\nADD file.tar.gz /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseCopyInsteadOfAdd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUseCopyInsteadOfAddJSON handles JSON-array form.
func TestIntegrationUseCopyInsteadOfAddJSON(t *testing.T) {
	src := "FROM alpine\nADD [\"file\",\"/dest\"]\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseCopyInsteadOfAdd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationUseCopyInsteadOfAddMultiple ensures any local source triggers.
func TestIntegrationUseCopyInsteadOfAddMultiple(t *testing.T) {
	src := "FROM alpine\nADD file file.tar.gz /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseCopyInsteadOfAdd()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationUseCopyInsteadOfAddNilDocument ensures nil input is handled.
func TestIntegrationUseCopyInsteadOfAddNilDocument(t *testing.T) {
	r := NewUseCopyInsteadOfAdd()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
