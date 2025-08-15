// file: cmd/docker-lint/errors_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
)

// TestRunConfigFlagMissingValue verifies that an error is returned when -c lacks a value.
func TestRunConfigFlagMissingValue(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"-c"}, &out, io.Discard, false)
	if err == nil || !strings.Contains(err.Error(), "missing config file") {
		t.Fatalf("expected missing config error, got %v", err)
	}
}

// TestRunConfigFileNotFound verifies that a missing config file causes an error.
func TestRunConfigFileNotFound(t *testing.T) {
	df := testDataPath("Dockerfile.good")
	var out bytes.Buffer
	err := run([]string{"-c", filepath.Join(t.TempDir(), "nope.yaml"), df}, &out, io.Discard, false)
	if err == nil {
		t.Fatalf("expected error for missing config file")
	}
}

// TestRunInvalidDefaultConfig ensures an invalid .docker-lint.yaml triggers an error.
func TestRunInvalidDefaultConfig(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".docker-lint.yaml")
	if err := os.WriteFile(cfgPath, []byte("::invalid"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	df := testDataPath("Dockerfile.good")
	t.Cleanup(func() { os.Remove(cfgPath) })
	t.Chdir(tmp)
	var out bytes.Buffer
	if err := run([]string{df}, &out, io.Discard, false); err == nil {
		t.Fatalf("expected error for invalid config")
	}
}

// TestExpandPathsInvalidPattern verifies that invalid glob patterns return an error.
func TestExpandPathsInvalidPattern(t *testing.T) {
	if _, err := expandPaths([]string{"["}); err == nil {
		t.Fatalf("expected glob error")
	}
}

// TestLintFileOpenError verifies that lintFile reports errors when files cannot be opened.
func TestLintFileOpenError(t *testing.T) {
	reg := engine.NewRegistry()
	if _, err := lintFile(context.Background(), reg, "does-not-exist"); err == nil {
		t.Fatalf("expected open error")
	}
}

// TestMainExitOnError ensures main exits with status 1 on failure.
func TestMainExitOnError(t *testing.T) {
	if os.Getenv("DOCKER_LINT_CRASHER") == "1" {
		os.Args = []string{"docker-lint"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainExitOnError")
	cmd.Env = append(os.Environ(), "DOCKER_LINT_CRASHER=1")
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() == 0 {
		t.Fatalf("expected exit code 1, got %v", err)
	}
}
