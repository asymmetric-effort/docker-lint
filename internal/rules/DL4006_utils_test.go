// file: internal/rules/DL4006_utils_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// TestIsNonPosixShell exercises detection of non-POSIX shells.
func TestIsNonPosixShell(t *testing.T) {
	shells := []string{"pwsh", "powershell"}
	if isNonPosixShell(nil, shells) {
		t.Fatalf("nil node should be false")
	}
	n := &parser.Node{Next: &parser.Node{Value: "pwsh"}}
	if !isNonPosixShell(n, shells) {
		t.Fatalf("expected pwsh to be non-posix")
	}
	posix := &parser.Node{Next: &parser.Node{Value: "/bin/bash"}}
	if isNonPosixShell(posix, shells) {
		t.Fatalf("unexpected non-posix result")
	}
}

// TestHasPipefailOption verifies detection of -o pipefail.
func TestHasPipefailOption(t *testing.T) {
	valid := map[string]bool{"/bin/bash": true}
	if hasPipefailOption(nil, valid) {
		t.Fatalf("nil node should be false")
	}
	n := &parser.Node{Next: &parser.Node{Value: "/bin/bash", Next: &parser.Node{Value: "-o", Next: &parser.Node{Value: "pipefail"}}}}
	if !hasPipefailOption(n, valid) {
		t.Fatalf("expected detection of pipefail option")
	}
	no := &parser.Node{Next: &parser.Node{Value: "/bin/bash"}}
	if hasPipefailOption(no, valid) {
		t.Fatalf("unexpected pipefail detection")
	}
	bad := &parser.Node{Next: &parser.Node{Value: "/bin/sh"}}
	if hasPipefailOption(bad, valid) {
		t.Fatalf("invalid shell should be false")
	}
}

// TestRunHasPipe covers shell and JSON RUN forms.
func TestRunHasPipe(t *testing.T) {
	if runHasPipe(nil) {
		t.Fatalf("nil node should be false")
	}
	jsonPipe := &parser.Node{Attributes: map[string]bool{"json": true}, Next: &parser.Node{Value: "echo", Next: &parser.Node{Value: "|"}}}
	if !runHasPipe(jsonPipe) {
		t.Fatalf("expected pipe in json form")
	}
	jsonNo := &parser.Node{Attributes: map[string]bool{"json": true}, Next: &parser.Node{Value: "echo"}}
	if runHasPipe(jsonNo) {
		t.Fatalf("unexpected pipe detection")
	}
	shPipe := &parser.Node{Next: &parser.Node{Value: "echo hi | grep h"}}
	if !runHasPipe(shPipe) {
		t.Fatalf("expected pipe in shell form")
	}
	shNo := &parser.Node{Next: &parser.Node{Value: "echo hi"}}
	if runHasPipe(shNo) {
		t.Fatalf("unexpected pipe detection")
	}
}
