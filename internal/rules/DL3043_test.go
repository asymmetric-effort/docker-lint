// file: internal/rules/DL3043_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationRequireOSVersionTagID validates rule identity.
func TestIntegrationRequireOSVersionTagID(t *testing.T) {
	if NewRequireOSVersionTag().ID() != "DL3043" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationRequireOSVersionTagViolations detects missing or floating OS tags.
func TestIntegrationRequireOSVersionTagViolations(t *testing.T) {
	r := NewRequireOSVersionTag()
	doc := &ir.Document{Stages: []*ir.Stage{
		{Index: 0, From: "ubuntu", Node: &parser.Node{StartLine: 1}},
		{Index: 1, From: "debian:stable", Node: &parser.Node{StartLine: 2}},
		{Index: 2, From: "alpine:3.19", Node: &parser.Node{StartLine: 3}},
	}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected two findings, got %d", len(findings))
	}
}

// TestIntegrationRequireOSVersionTagClean ensures compliant images pass.
func TestIntegrationRequireOSVersionTagClean(t *testing.T) {
	r := NewRequireOSVersionTag()
	doc := &ir.Document{Stages: []*ir.Stage{
		{Index: 0, From: "ubuntu:22.04", Node: &parser.Node{StartLine: 1}},
		{Index: 1, From: "golang:1.21", Node: &parser.Node{StartLine: 2}},
	}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireOSVersionTagAlias ensures alias references are exempt.
func TestIntegrationRequireOSVersionTagAlias(t *testing.T) {
	r := NewRequireOSVersionTag()
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

// TestIntegrationRequireOSVersionTagScratchAndVariable allows scratch and variable images.
func TestIntegrationRequireOSVersionTagScratchAndVariable(t *testing.T) {
	r := NewRequireOSVersionTag()
	doc := &ir.Document{Stages: []*ir.Stage{
		{Index: 0, From: "scratch", Node: &parser.Node{StartLine: 1}},
		{Index: 1, From: "$BASE", Node: &parser.Node{StartLine: 2}},
	}}
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationRequireOSVersionTagNilDocument ensures graceful handling of nil input.
func TestIntegrationRequireOSVersionTagNilDocument(t *testing.T) {
	r := NewRequireOSVersionTag()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
}
