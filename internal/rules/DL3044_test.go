// file: internal/rules/DL3044_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationDnfVersionPinID validates rule identity.
func TestIntegrationDnfVersionPinID(t *testing.T) {
	if NewDnfVersionPin().ID() != "DL3044" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationDnfVersionPinViolation detects unpinned dnf installs.
func TestIntegrationDnfVersionPinViolation(t *testing.T) {
	src := "FROM fedora\nRUN dnf install -y curl\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfVersionPin()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationMicrodnfVersionPinViolation detects unpinned microdnf installs.
func TestIntegrationMicrodnfVersionPinViolation(t *testing.T) {
	src := "FROM registry.access.redhat.com/ubi9/ubi\nRUN microdnf install -y python3 make\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfVersionPin()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationDnfVersionPinClean ensures pinned installs pass.
func TestIntegrationDnfVersionPinClean(t *testing.T) {
	src := "FROM fedora\nRUN dnf install -y curl-7.76.1-5.fc34 make-4.3-3.fc34\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfVersionPin()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationDnfVersionPinVariable ensures variable pins are allowed.
func TestIntegrationDnfVersionPinVariable(t *testing.T) {
	src := "FROM fedora\nARG CURL_VER=7.76.1-5.fc34\nRUN dnf install -y curl-${CURL_VER}\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build doc: %v", err)
	}
	r := NewDnfVersionPin()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationDnfVersionPinNil ensures graceful handling of nil input.
func TestIntegrationDnfVersionPinNil(t *testing.T) {
	r := NewDnfVersionPin()
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
	if f, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", f, err)
	}
}
