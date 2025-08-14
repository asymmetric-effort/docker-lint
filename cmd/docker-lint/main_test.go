// file: cmd/docker-lint/main_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
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
