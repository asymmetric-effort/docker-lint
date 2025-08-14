// file: internal/rules/no_latest_tag_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

func TestNoLatestTagID(t *testing.T) {
	if NewNoLatestTag().ID() != "DL3007" {
		t.Fatalf("unexpected id")
	}
}

func TestNoLatestTagViolation(t *testing.T) {
	r := NewNoLatestTag()
	doc := &ir.Document{Stages: []*ir.Stage{
		{Index: 0, From: "alpine", Node: &parser.Node{StartLine: 1}},
		{Index: 1, From: "alpine:latest", Node: &parser.Node{StartLine: 2}},
	}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

func TestNoLatestTagClean(t *testing.T) {
	r := NewNoLatestTag()
	doc := &ir.Document{Stages: []*ir.Stage{{Index: 0, From: "alpine:3.19", Node: &parser.Node{StartLine: 1}}, {Index: 1, From: "ubuntu@sha256:deadbeef", Node: &parser.Node{StartLine: 2}}}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}
