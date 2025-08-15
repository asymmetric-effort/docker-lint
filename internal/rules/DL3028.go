// file: internal/rules/DL3028.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// pinGemVersions enforces version pinning in gem install commands.
type pinGemVersions struct{}

// NewPinGemVersions constructs the rule.
func NewPinGemVersions() engine.Rule { return pinGemVersions{} }

// ID returns the rule identifier.
func (pinGemVersions) ID() string { return "DL3028" }

// Check inspects RUN instructions for unpinned gem installs.
func (pinGemVersions) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := splitRunSegments(n)
		for _, seg := range segments {
			if violatesGemPin(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3028",
					Message: "Pin versions in gem install. Instead of `gem install <gem>` use `gem install <gem>:<version>`",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// violatesGemPin reports whether a command segment uses gem install without version.
func violatesGemPin(seg []string) bool {
	for i := 0; i < len(seg); i++ {
		if seg[i] != "gem" {
			continue
		}
		if i+1 >= len(seg) || (seg[i+1] != "install" && seg[i+1] != "i") {
			continue
		}
		args := seg[i+2:]
		for _, a := range args {
			if a == "-v" || a == "--version" || strings.HasPrefix(a, "--version=") {
				return false
			}
		}
		for _, a := range args {
			if a == "--" {
				break
			}
			if strings.HasPrefix(a, "-") {
				continue
			}
			if a != "install" && a != "i" && !strings.Contains(a, ":") {
				return true
			}
		}
	}
	return false
}
