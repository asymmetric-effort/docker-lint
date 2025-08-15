// file: internal/rules/DL3006_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationRequireTagID validates rule identity.
func TestIntegrationRequireTagID(t *testing.T) {
	if NewRequireTag().ID() != "DL3006" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationRequireTagViolation detects untagged images.
func TestIntegrationRequireTagViolation(t *testing.T) {
	r := NewRequireTag()
	doc := &ir.Document{Stages: []*ir.Stage{{Index: 0, From: "alpine", Node: &parser.Node{StartLine: 1}}}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationRequireTagAlias ensures alias references are exempt.
func TestIntegrationRequireTagAlias(t *testing.T) {
	r := NewRequireTag()
	doc := &ir.Document{Stages: []*ir.Stage{
		{Index: 0, From: "alpine:3.19", Name: "base", Node: &parser.Node{StartLine: 1}},
		{Index: 1, From: "base", Node: &parser.Node{StartLine: 2}},
	}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireTagClean ensures compliant images pass.
func TestIntegrationRequireTagClean(t *testing.T) {
	r := NewRequireTag()
	doc := &ir.Document{Stages: []*ir.Stage{
		{Index: 0, From: "alpine:3.19", Node: &parser.Node{StartLine: 1}},
		{Index: 1, From: "ubuntu@sha256:deadbeef", Node: &parser.Node{StartLine: 2}},
	}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireTagNilDocument ensures graceful handling of nil input.
func TestIntegrationRequireTagNilDocument(t *testing.T) {
	r := NewRequireTag()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
}
