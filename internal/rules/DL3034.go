// file: internal/rules/DL3034.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireZypperYes enforces non-interactive zypper commands.
type requireZypperYes struct{}

// NewRequireZypperYes constructs the rule.
func NewRequireZypperYes() engine.Rule { return requireZypperYes{} }

// ID returns the rule identifier.
func (requireZypperYes) ID() string { return "DL3034" }

// Check verifies zypper commands include a non-interactive switch.
func (requireZypperYes) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if isZypperAction(seg) && !hasZypperYes(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3034",
					Message: "Non-interactive switch missing from `zypper` command: `zypper install -y`",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isZypperAction reports if command is zypper install/remove/etc.
func isZypperAction(seg []string) bool {
	if len(seg) < 2 || seg[0] != "zypper" {
		return false
	}
	actions := []string{"install", "in", "remove", "rm", "source-install", "si", "patch"}
	for _, a := range actions {
		if seg[1] == a {
			return true
		}
	}
	return false
}

// hasZypperYes detects non-interactive flags.
func hasZypperYes(seg []string) bool {
	for _, t := range seg {
		if t == "-y" || t == "-n" || t == "--non-interactive" || t == "--no-confirm" {
			return true
		}
	}
	return false
}
