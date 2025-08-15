// file: internal/rules/run_commands_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"reflect"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// TestExtractCommands exercises extractCommands across JSON and shell forms.
func TestExtractCommands(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if extractCommands(nil) != nil {
			t.Fatalf("expected nil for nil node")
		}
	})
	t.Run("json", func(t *testing.T) {
		n := &parser.Node{Attributes: map[string]bool{"json": true}, Next: &parser.Node{Value: "CMD"}}
		got := extractCommands(n)
		want := []string{"cmd"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("shell", func(t *testing.T) {
		n := &parser.Node{Next: &parser.Node{Value: "echo hi && ls"}}
		got := extractCommands(n)
		want := []string{"echo", "ls"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("bad shell", func(t *testing.T) {
		n := &parser.Node{Next: &parser.Node{Value: "echo 'unterminated"}}
		if extractCommands(n) != nil {
			t.Fatalf("expected nil on parse error")
		}
	})
}

// TestCommandNames verifies command extraction with shell connectors.
func TestCommandNames(t *testing.T) {
	tokens := []string{"echo", "hi", "&&", "ls", "||", "cat", "|", "grep", ";", "sed"}
	want := []string{"echo", "ls", "cat", "grep", "sed"}
	if got := commandNames(tokens); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

// TestLowerSegments confirms all segments are lowercased.
func TestLowerSegments(t *testing.T) {
	segs := [][]string{{"Echo", "Hi"}, {"LS"}}
	got := lowerSegments(segs)
	want := [][]string{{"echo", "hi"}, {"ls"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
