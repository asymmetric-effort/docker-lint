// file: internal/rules/stage_utils_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// TestCopyFromFlag verifies extraction of the --from flag.
func TestCopyFromFlag(t *testing.T) {
	n := &parser.Node{Flags: []string{"--from=builder", "--chown=0"}}
	v, ok := copyFromFlag(n)
	if !ok || v != "builder" {
		t.Fatalf("expected builder, got %q %v", v, ok)
	}
	n = &parser.Node{Flags: []string{"--chown=0"}}
	if _, ok := copyFromFlag(n); ok {
		t.Fatalf("expected no flag")
	}
}

// TestStageAlias verifies detection of stage aliases in FROM instructions.
func TestStageAlias(t *testing.T) {
	n := &parser.Node{Value: "from"}
	n.Next = &parser.Node{Value: "alpine", Next: &parser.Node{Value: "as", Next: &parser.Node{Value: "builder"}}}
	if a := stageAlias(n); a != "builder" {
		t.Fatalf("expected builder, got %q", a)
	}
	if a := stageAlias(&parser.Node{Value: "from", Next: &parser.Node{Value: "alpine"}}); a != "" {
		t.Fatalf("expected empty alias, got %q", a)
	}
	if a := stageAlias(nil); a != "" {
		t.Fatalf("expected empty alias for nil input")
	}
}
