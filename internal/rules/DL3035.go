// file: internal/rules/DL3035.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// forbidZypperDistUpgrade prohibits zypper dist-upgrade usage.
type forbidZypperDistUpgrade struct{}

// NewForbidZypperDistUpgrade constructs the rule.
func NewForbidZypperDistUpgrade() engine.Rule { return forbidZypperDistUpgrade{} }

// ID returns the rule identifier.
func (forbidZypperDistUpgrade) ID() string { return "DL3035" }

// Check flags zypper dist-upgrade commands.
func (forbidZypperDistUpgrade) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if isZypperDistUpgrade(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3035",
					Message: "Do not use `zypper dist-upgrade`.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isZypperDistUpgrade reports if segment runs zypper dist-upgrade/dup.
func isZypperDistUpgrade(seg []string) bool {
	if len(seg) < 2 || seg[0] != "zypper" {
		return false
	}
	if seg[1] == "dist-upgrade" || seg[1] == "dup" {
		return true
	}
	return false
}
