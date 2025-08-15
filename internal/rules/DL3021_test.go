// file: internal/rules/DL3021_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationCopyDestSlashID validates rule identity.
func TestIntegrationCopyDestSlashID(t *testing.T) {
	if NewCopyDestEndsWithSlash().ID() != "DL3021" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationCopyDestSlashViolation detects missing trailing slash with multiple sources.
func TestIntegrationCopyDestSlashViolation(t *testing.T) {
	src := "FROM alpine\nCOPY file1 file2 /opt\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewCopyDestEndsWithSlash()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationCopyDestSlashClean ensures compliant Dockerfiles pass.
func TestIntegrationCopyDestSlashClean(t *testing.T) {
	src := "FROM alpine\nCOPY file1 file2 /opt/\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewCopyDestEndsWithSlash()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyDestSlashSingleSource verifies single-source COPY is ignored.
func TestIntegrationCopyDestSlashSingleSource(t *testing.T) {
	src := "FROM alpine\nCOPY file1 /opt\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewCopyDestEndsWithSlash()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyDestSlashNilDocument ensures graceful handling of nil input.
func TestIntegrationCopyDestSlashNilDocument(t *testing.T) {
	r := NewCopyDestEndsWithSlash()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
