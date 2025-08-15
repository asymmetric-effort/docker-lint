// file: internal/rules/ruleutil_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"reflect"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// TestSplitRunSegments covers JSON, shell, and error scenarios.
func TestSplitRunSegments(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if splitRunSegments(nil) != nil {
			t.Fatalf("expected nil for nil node")
		}
	})
	t.Run("json", func(t *testing.T) {
		n := &parser.Node{Attributes: map[string]bool{"json": true}, Next: &parser.Node{Value: "CMD"}}
		got := splitRunSegments(n)
		want := [][]string{{"cmd"}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("shell", func(t *testing.T) {
		n := &parser.Node{Next: &parser.Node{Value: "echo hi && ls"}}
		got := splitRunSegments(n)
		want := [][]string{{"echo", "hi"}, {"ls"}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("bad shell", func(t *testing.T) {
		n := &parser.Node{Next: &parser.Node{Value: "echo 'unterminated"}}
		if splitRunSegments(n) != nil {
			t.Fatalf("expected nil on parse error")
		}
	})
}

// TestLowerSlice ensures lowerSlice returns a lowercase copy.
func TestLowerSlice(t *testing.T) {
	got := lowerSlice([]string{"A", "b"})
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
