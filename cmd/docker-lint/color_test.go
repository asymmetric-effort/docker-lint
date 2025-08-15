// file: cmd/docker-lint/color_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/sam-caldwell/ansi"
)

// TestRunColoredWarnings verifies that warnings are printed in yellow when color is enabled.
func TestRunColoredWarnings(t *testing.T) {
	df := testDataPath("Dockerfile.bad")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	if err := run([]string{df}, &out, &errBuf, true); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(errBuf.String(), ansi.CodeFgYellow) {
		t.Fatalf("expected yellow output, got %q", errBuf.String())
	}
}

// TestRunColoredSuccess verifies that success messages are printed in green when color is enabled.
func TestRunColoredSuccess(t *testing.T) {
	df := testDataPath("Dockerfile.good")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	if err := run([]string{df}, &out, &errBuf, true); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(errBuf.String(), ansi.CodeFgGreen) {
		t.Fatalf("expected green output, got %q", errBuf.String())
	}
}

// TestRunNoColor verifies that color codes are absent when color is disabled.
func TestRunNoColor(t *testing.T) {
	df := testDataPath("Dockerfile.bad")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	if err := run([]string{df}, &out, &errBuf, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if strings.Contains(errBuf.String(), ansi.CodeFgYellow) || strings.Contains(errBuf.String(), ansi.CodeFgGreen) {
		t.Fatalf("unexpected color codes in output: %q", errBuf.String())
	}
}

// TestPrintErrorColor verifies that errors are printed in red when color is enabled.
func TestPrintErrorColor(t *testing.T) {
	var errBuf bytes.Buffer
	printError(&errBuf, true, fmt.Errorf("boom"))
	if !strings.Contains(errBuf.String(), ansi.CodeFgRed) {
		t.Fatalf("expected red output, got %q", errBuf.String())
	}
}

// TestPrintErrorNoColor verifies that errors are not colorized when color is disabled.
func TestPrintErrorNoColor(t *testing.T) {
	var errBuf bytes.Buffer
	printError(&errBuf, false, fmt.Errorf("boom"))
	if strings.Contains(errBuf.String(), ansi.CodeFgRed) {
		t.Fatalf("unexpected color codes in output: %q", errBuf.String())
	}
}
