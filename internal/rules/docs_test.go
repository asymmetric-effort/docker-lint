// file: internal/rules/docs_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// ruleAndDocDirs returns the absolute rule and documentation directories.
func ruleAndDocDirs(t *testing.T) (string, string) {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..")
	ruleDir := filepath.Join(root, "internal", "rules")
	docDir := filepath.Join(root, "docs", "rules")
	return ruleDir, docDir
}

// TestRuleDocsExist ensures every rule has corresponding documentation.
func TestRuleDocsExist(t *testing.T) {
	ruleDir, docDir := ruleAndDocDirs(t)
	entries, err := os.ReadDir(ruleDir)
	if err != nil {
		t.Fatalf("read rules: %v", err)
	}
	re := regexp.MustCompile(`^DL\d{4}\.go$`)
	for _, e := range entries {
		name := e.Name()
		if !re.MatchString(name) {
			continue
		}
		doc := filepath.Join(docDir, strings.TrimSuffix(name, ".go")+".md")
		if _, err := os.Stat(doc); err != nil {
			t.Errorf("missing documentation for %s", name)
		}
	}
}

// TestDocsMatchRules ensures each documentation file corresponds to a rule and
// starts with the rule identifier.
func TestDocsMatchRules(t *testing.T) {
	ruleDir, docDir := ruleAndDocDirs(t)
	entries, err := os.ReadDir(docDir)
	if err != nil {
		t.Fatalf("read docs: %v", err)
	}
	re := regexp.MustCompile(`^DL\d{4}\.md$`)
	for _, e := range entries {
		name := e.Name()
		if name == "README.md" || !re.MatchString(name) {
			continue
		}
		rule := filepath.Join(ruleDir, strings.TrimSuffix(name, ".md")+".go")
		if _, err := os.Stat(rule); err != nil {
			t.Errorf("documentation exists for missing rule %s", name)
			continue
		}
		path := filepath.Join(docDir, name)
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("open %s: %v", path, err)
		}
		scanner := bufio.NewScanner(f)
		if !scanner.Scan() {
			t.Errorf("documentation %s is empty", name)
			f.Close()
			continue
		}
		title := scanner.Text()
		f.Close()
		id := strings.TrimSuffix(name, ".md")
		if !strings.HasPrefix(title, "# "+id) {
			t.Errorf("documentation %s has incorrect title %q", name, title)
		}
	}
}
