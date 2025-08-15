// file: cmd/docker-lint/main_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/sam-caldwell/ansi"
)

func TestIntegrationMain(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"docker-lint", testDataPath("Dockerfile.bad")}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout
	out, _ := io.ReadAll(r)

	var findings []engine.Finding
	if err := json.Unmarshal(out, &findings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(findings) == 0 {
		t.Fatalf("expected findings")
	}
}

// TestMainNoColorFlag verifies that the --no-color flag disables colored output.
func TestMainNoColorFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"docker-lint", "--no-color", testDataPath("Dockerfile.good")}

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w

	main()

	w.Close()
	os.Stderr = oldStderr
	out, _ := io.ReadAll(r)
	if strings.Contains(string(out), ansi.CodeFgGreen) {
		t.Fatalf("unexpected color codes in stderr: %q", out)
	}
}
