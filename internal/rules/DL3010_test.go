// file: internal/rules/DL3010_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationUseADDForArchivesID validates rule identity.
func TestIntegrationUseADDForArchivesID(t *testing.T) {
	if NewUseADDForArchives().ID() != "DL3010" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationUseADDForArchivesViolation detects COPY of archives into directories.
func TestIntegrationUseADDForArchivesViolation(t *testing.T) {
	src := "FROM alpine\nCOPY app.tar.gz /opt/\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseADDForArchives()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationUseADDForArchivesFileDest ignores archives copied to file paths.
func TestIntegrationUseADDForArchivesFileDest(t *testing.T) {
	src := "FROM alpine\nCOPY app.tar.gz /opt/app.tar.gz\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseADDForArchives()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUseADDForArchivesClean ensures compliant Dockerfiles pass.
func TestIntegrationUseADDForArchivesClean(t *testing.T) {
	src := "FROM alpine\nCOPY app.txt /opt/\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseADDForArchives()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUseADDForArchivesMultiStage ignores --from copies.
func TestIntegrationUseADDForArchivesMultiStage(t *testing.T) {
	src := "FROM scratch AS build\nFROM scratch\nCOPY --from=build app.tar.gz /opt/\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewUseADDForArchives()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationUseADDForArchivesNilDocument ensures graceful handling of nil input.
func TestIntegrationUseADDForArchivesNilDocument(t *testing.T) {
	r := NewUseADDForArchives()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
