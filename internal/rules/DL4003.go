// file: internal/rules/DL4003.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// singleCmd enforces a single CMD instruction per build stage.
type singleCmd struct{}

// NewSingleCmd constructs the rule.
func NewSingleCmd() engine.Rule { return singleCmd{} }

// ID returns the rule identifier.
func (singleCmd) ID() string { return "DL4003" }

// Check scans each stage for multiple CMD instructions.
func (singleCmd) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	seen := false
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "from") {
			seen = false
			continue
		}
		if strings.EqualFold(n.Value, "cmd") {
			if seen {
				findings = append(findings, engine.Finding{
					RuleID:  "DL4003",
					Message: "Multiple CMD instructions found. If you list more than one CMD then only the last CMD will take effect",
					Line:    n.StartLine,
				})
				continue
			}
			seen = true
		}
	}
	return findings, nil
}
