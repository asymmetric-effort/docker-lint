// file: internal/ir/document_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package ir

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func TestBuildDocument(t *testing.T) {
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
