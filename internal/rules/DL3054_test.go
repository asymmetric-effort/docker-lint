// file: internal/rules/DL3054_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
        "context"
        "strings"
        "testing"

        "github.com/moby/buildkit/frontend/dockerfile/parser"

        "github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationLabelSPDXValidID validates rule identity.
func TestIntegrationLabelSPDXValidID(t *testing.T) {
        if NewLabelSPDXValid(nil).ID() != "DL3054" {
                t.Fatalf("unexpected id")
        }
}

// TestIntegrationLabelSPDXValidViolation detects invalid SPDX license identifiers.
func TestIntegrationLabelSPDXValidViolation(t *testing.T) {
	src := "FROM scratch\nLABEL license=not@spdx\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"license": LabelTypeSPDX}
	r := NewLabelSPDXValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
        if len(findings) != 1 {
                t.Fatalf("expected one finding, got %d", len(findings))
        }
}

// TestIntegrationLabelSPDXValidClean ensures valid SPDX identifiers pass.
func TestIntegrationLabelSPDXValidClean(t *testing.T) {
	src := "FROM scratch\nLABEL license=MIT\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"license": LabelTypeSPDX}
	r := NewLabelSPDXValid(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
        if len(findings) != 0 {
                t.Fatalf("expected no findings, got %d", len(findings))
        }
}

// TestIntegrationLabelSPDXValidNilDocument ensures nil documents are handled gracefully.
func TestIntegrationLabelSPDXValidNilDocument(t *testing.T) {
	r := NewLabelSPDXValid(nil)
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
