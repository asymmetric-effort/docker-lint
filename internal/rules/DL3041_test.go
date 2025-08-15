// file: internal/rules/DL3041_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationDnfNoUpgradeID validates rule identity.
func TestIntegrationDnfNoUpgradeID(t *testing.T) {
	if NewDnfNoUpgrade().ID() != "DL3041" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationDnfNoUpgradeViolation detects dnf upgrade usage.
func TestIntegrationDnfNoUpgradeViolation(t *testing.T) {
	src := "FROM quay.io/ubi9/ubi\nRUN dnf -y upgrade\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfNoUpgrade()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationMicrodnfNoUpgradeViolation detects microdnf update usage.
func TestIntegrationMicrodnfNoUpgradeViolation(t *testing.T) {
	src := "FROM registry.access.redhat.com/ubi9/ubi\nRUN microdnf update -y\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfNoUpgrade()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationDnfNoUpgradeClean ensures compliant Dockerfiles pass.
func TestIntegrationDnfNoUpgradeClean(t *testing.T) {
	src := "FROM fedora\nRUN dnf install -y httpd\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfNoUpgrade()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationDnfNoUpgradeNil ensures graceful handling of nil input.
func TestIntegrationDnfNoUpgradeNil(t *testing.T) {
	r := NewDnfNoUpgrade()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
