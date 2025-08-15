// file: internal/rules/DL3047_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationApkCacheCleanupID validates rule identity.
func TestIntegrationApkCacheCleanupID(t *testing.T) {
	if NewApkCacheCleanup().ID() != "DL3047" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationApkCacheCleanupViolation reports missing cleanup without no-cache.
func TestIntegrationApkCacheCleanupViolation(t *testing.T) {
	src := "FROM alpine\nRUN apk add curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewApkCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationApkCacheCleanupNoCache passes with --no-cache flag.
func TestIntegrationApkCacheCleanupNoCache(t *testing.T) {
	src := "FROM alpine\nRUN apk add --no-cache curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewApkCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationApkCacheCleanupRemove removes cache manually.
func TestIntegrationApkCacheCleanupRemove(t *testing.T) {
	src := "FROM alpine\nRUN apk add curl && rm -rf /var/cache/apk/*\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewApkCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationApkCacheCleanupFindDelete uses find -delete.
func TestIntegrationApkCacheCleanupFindDelete(t *testing.T) {
	src := "FROM alpine\nRUN apk add curl && find /var/cache/apk -type f -delete\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewApkCacheCleanup()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationApkCacheCleanupNil ensures nil documents handled.
func TestIntegrationApkCacheCleanupNil(t *testing.T) {
	r := NewApkCacheCleanup()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
