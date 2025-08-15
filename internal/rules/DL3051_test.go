// file: internal/rules/DL3051_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
        "context"
        "strings"
        "testing"

        "github.com/moby/buildkit/frontend/dockerfile/parser"

        "github.com/asymmetric-effort/docker-lint/internal/ir"
)

// TestIntegrationLabelNotEmptyID validates rule identity.
func TestIntegrationLabelNotEmptyID(t *testing.T) {
        if NewLabelNotEmpty(nil).ID() != "DL3051" {
                t.Fatalf("unexpected id")
        }
}

// TestIntegrationLabelNotEmptyViolation detects empty label values.
func TestIntegrationLabelNotEmptyViolation(t *testing.T) {
	src := "FROM scratch\nLABEL foo=\"\"\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"foo": LabelTypeString}
	r := NewLabelNotEmpty(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
        if len(findings) != 1 {
                t.Fatalf("expected one finding, got %d", len(findings))
        }
}

// TestIntegrationLabelNotEmptyClean ensures populated values pass.
func TestIntegrationLabelNotEmptyClean(t *testing.T) {
	src := "FROM scratch\nLABEL foo=bar\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"foo": LabelTypeString}
	r := NewLabelNotEmpty(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
        if len(findings) != 0 {
                t.Fatalf("expected no findings, got %d", len(findings))
        }
}

// TestIntegrationLabelNotEmptyWhitespace flags values containing only whitespace.
func TestIntegrationLabelNotEmptyWhitespace(t *testing.T) {
	src := "FROM scratch\nLABEL foo=\" \"\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"foo": LabelTypeString}
	r := NewLabelNotEmpty(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
        if len(findings) != 1 {
                t.Fatalf("expected one finding, got %d", len(findings))
        }
}

// TestIntegrationLabelNotEmptyOverwrite verifies later non-empty labels override earlier empty ones.
func TestIntegrationLabelNotEmptyOverwrite(t *testing.T) {
	src := "FROM scratch\nLABEL foo=\"\"\nLABEL foo=\"bar\"\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	schema := LabelSchema{"foo": LabelTypeString}
	r := NewLabelNotEmpty(schema)
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}

	src = "FROM scratch\nLABEL foo=\"bar\"\nLABEL foo=\"\"\n"
	res, err = parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err = ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	findings, err = r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
        if len(findings) != 1 {
                t.Fatalf("expected one finding, got %d", len(findings))
        }
}

// TestIntegrationLabelNotEmptyNilDocument ensures nil documents are handled gracefully.
func TestIntegrationLabelNotEmptyNilDocument(t *testing.T) {
	r := NewLabelNotEmpty(nil)
	if f, err := r.Check(context.Background(), nil); err != nil || len(f) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", f, err)
	}
}
