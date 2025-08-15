// file: internal/engine/ignore_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package engine

import (
	"reflect"
	"testing"
)

// TestParseIgnorePragmaUnicode ensures multi-byte characters before the directive do not break parsing.
func TestParseIgnorePragmaUnicode(t *testing.T) {
	input := "# İ hadolint ignore=DL3007, DL2000"
	got := parseIgnorePragma(input)
	want := []string{"dl3007", "dl2000"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseIgnorePragma(%q) = %v; want %v", input, got, want)
	}
}
