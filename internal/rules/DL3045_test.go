// file: internal/rules/DL3045_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationCopyFromExternalDigestID validates rule identity.
func TestIntegrationCopyFromExternalDigestID(t *testing.T) {
	if NewCopyFromExternalDigest().ID() != "DL3045" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationCopyFromExternalDigestViolation detects missing digests.
func TestIntegrationCopyFromExternalDigestViolation(t *testing.T) {
	src := "FROM alpine\nCOPY --from=ubuntu:22.04 /bin/bash /bin/bash\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromExternalDigest()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationCopyFromExternalDigestCompliant ensures digest-pinned images pass.
func TestIntegrationCopyFromExternalDigestCompliant(t *testing.T) {
	src := "FROM alpine\nCOPY --from=ubuntu:22.04@sha256:deadbeef /bin/bash /bin/bash\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromExternalDigest()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyFromExternalDigestStageAlias ignores intra-file stages.
func TestIntegrationCopyFromExternalDigestStageAlias(t *testing.T) {
	src := "FROM alpine AS build\nCOPY --from=build /src /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromExternalDigest()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyFromExternalDigestScratch skips scratch source.
func TestIntegrationCopyFromExternalDigestScratch(t *testing.T) {
	src := "FROM alpine\nCOPY --from=scratch /bin/bash /bin/bash\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromExternalDigest()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyFromExternalDigestNumeric skips numeric stage references.
func TestIntegrationCopyFromExternalDigestNumeric(t *testing.T) {
	src := "FROM alpine\nFROM alpine AS build\nCOPY --from=0 /src /dest\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewCopyFromExternalDigest()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationCopyFromExternalDigestNil ensures graceful handling of nil input.
func TestIntegrationCopyFromExternalDigestNil(t *testing.T) {
	r := NewCopyFromExternalDigest()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
