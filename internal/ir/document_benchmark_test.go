// file: internal/ir/document_benchmark_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package ir

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// BenchmarkBuildDocument measures the performance of building a Document
// from a parsed Dockerfile AST.
func BenchmarkBuildDocument(b *testing.B) {
	src := "FROM alpine AS base\nRUN echo hi\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		b.Fatalf("parse failed: %v", err)
	}
	ast := res.AST
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := BuildDocument("Dockerfile", ast); err != nil {
			b.Fatalf("build document: %v", err)
		}
	}
}
