// file: internal/ir/document_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package ir

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// TestIntegrationBuildDocument verifies basic document construction.
func TestIntegrationBuildDocument(t *testing.T) {
	src := "FROM alpine:3.19 AS base\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	if len(doc.Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(doc.Stages))
	}
	st := doc.Stages[0]
	if st.From != "alpine:3.19" || st.Name != "base" {
		t.Fatalf("unexpected stage: %+v", st)
	}
}

// TestIntegrationBuildDocumentRetainsAST verifies the original AST is stored.
func TestIntegrationBuildDocumentRetainsAST(t *testing.T) {
	src := "FROM scratch\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	if doc.AST != res.AST {
		t.Fatalf("document AST mismatch")
	}
	if len(doc.Stages) != 1 || doc.Stages[0].Node != res.AST.Children[0] {
		t.Fatalf("stage node mismatch")
	}
}

// TestIntegrationBuildDocumentMultipleStages ensures indexing across multiple FROMs.
func TestIntegrationBuildDocumentMultipleStages(t *testing.T) {
	src := "FROM alpine AS base\nFROM scratch\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	if len(doc.Stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(doc.Stages))
	}
	first, second := doc.Stages[0], doc.Stages[1]
	if first.Index != 0 || first.From != "alpine" || first.Name != "base" {
		t.Fatalf("unexpected first stage: %+v", first)
	}
	if second.Index != 1 || second.From != "scratch" || second.Name != "" {
		t.Fatalf("unexpected second stage: %+v", second)
	}
}
