// file: internal/rules/DL3040_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationDnfCacheCleanupID validates rule identity.
func TestIntegrationDnfCacheCleanupID(t *testing.T) {
	if NewDnfCacheCleanup().ID() != "DL3040" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationDnfCacheCleanupViolation detects missing cleanup after dnf install.
func TestIntegrationDnfCacheCleanupViolation(t *testing.T) {
	src := "FROM fedora\nRUN dnf install -y httpd\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDnfCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationDnfCacheCleanupClean ensures cleanup after install passes.
func TestIntegrationDnfCacheCleanupClean(t *testing.T) {
	src := "FROM fedora\nRUN dnf install -y httpd && dnf clean all\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDnfCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationMicrodnfCleanup verifies microdnf cleanup.
func TestIntegrationMicrodnfCleanup(t *testing.T) {
	src := "FROM fedora\nRUN microdnf install -y ca-certificates && microdnf clean all\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDnfCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationDnfCacheCleanupRM verifies explicit cache removal via rm.
func TestIntegrationDnfCacheCleanupRM(t *testing.T) {
	src := "FROM fedora\nRUN dnf install -y httpd && rm -rf /var/cache/dnf\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDnfCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationDnfCacheCleanupOrder ensures cleanup preceding install fails.
func TestIntegrationDnfCacheCleanupOrder(t *testing.T) {
	src := "FROM fedora\nRUN dnf clean all && dnf install -y httpd\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewDnfCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationDnfCacheCleanupNil ensures nil documents are handled.
func TestIntegrationDnfCacheCleanupNil(t *testing.T) {
	r := NewDnfCacheCleanup()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
