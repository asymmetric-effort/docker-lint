// file: internal/rules/docs_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// TestRuleDocsExist ensures every rule has corresponding documentation.
func TestRuleDocsExist(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..")
	ruleDir := filepath.Join(root, "internal", "rules")
	docDir := filepath.Join(root, "docs", "rules")
	entries, err := os.ReadDir(ruleDir)
	if err != nil {
		t.Fatalf("read rules: %v", err)
	}
	re := regexp.MustCompile(`^DL\d+\.go$`)
	for _, e := range entries {
		name := e.Name()
		if !re.MatchString(name) {
			continue
		}
		doc := filepath.Join(docDir, name[:len(name)-3]+".md")
		if _, err := os.Stat(doc); err != nil {
			t.Errorf("missing documentation for %s", name)
		}
	}
}
