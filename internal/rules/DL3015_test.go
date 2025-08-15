// file: internal/rules/DL3015_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationAptNoInstallRecommendsID validates rule identity.
func TestIntegrationAptNoInstallRecommendsID(t *testing.T) {
	if NewAptNoInstallRecommends().ID() != "DL3015" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationAptNoInstallRecommendsViolation detects missing no-install-recommends.
func TestIntegrationAptNoInstallRecommendsViolation(t *testing.T) {
	src := "FROM alpine\nRUN apt-get install -y gcc\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAptNoInstallRecommends()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

// TestIntegrationAptNoInstallRecommendsClean ensures compliant Dockerfiles pass.
func TestIntegrationAptNoInstallRecommendsClean(t *testing.T) {
	src := "FROM alpine\nRUN apt-get install --no-install-recommends -y gcc\nRUN apt-get -o APT::Install-Recommends=false install gcc\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewAptNoInstallRecommends()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationAptNoInstallRecommendsNilDocument ensures graceful handling of nil input.
func TestIntegrationAptNoInstallRecommendsNilDocument(t *testing.T) {
	r := NewAptNoInstallRecommends()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}
