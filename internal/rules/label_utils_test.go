// file: internal/rules/label_utils_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// parseLabelNode parses the provided LABEL instruction and returns its AST node.
func parseLabelNode(t *testing.T, instr string) *parser.Node {
	t.Helper()
	src := "FROM scratch\n" + instr + "\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(res.AST.Children) < 2 {
		t.Fatalf("expected label node, got %v", res.AST.Children)
	}
	return res.AST.Children[1]
}

// TestCollectLabelPairs exercises extraction of key-value pairs from LABEL instructions.
func TestCollectLabelPairs(t *testing.T) {
	t.Run("nil node", func(t *testing.T) {
		if pairs := collectLabelPairs(nil); len(pairs) != 0 {
			t.Fatalf("expected no pairs, got %d", len(pairs))
		}
	})

	t.Run("multiple pairs", func(t *testing.T) {
		ln := parseLabelNode(t, "LABEL foo=bar baz=qux")
		pairs := collectLabelPairs(ln)
		expected := []labelPair{{Key: "foo", Value: "bar"}, {Key: "baz", Value: "qux"}}
		if len(pairs) != len(expected) {
			t.Fatalf("expected %d pairs, got %d", len(expected), len(pairs))
		}
		for i, p := range expected {
			if pairs[i] != p {
				t.Fatalf("pair %d mismatch: expected %v, got %v", i, p, pairs[i])
			}
		}
	})
}

// TestInSchema verifies key lookup against a label schema.
func TestInSchema(t *testing.T) {
	schema := LabelSchema{"foo": LabelTypeString}
	if !inSchema(schema, "foo") {
		t.Fatalf("expected key to be in schema")
	}
	if inSchema(schema, "bar") {
		t.Fatalf("unexpected key found in schema")
	}
}
