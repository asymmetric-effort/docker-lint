// file: internal/rules/DL3027.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// noAptCommand warns against using apt instead of apt-get or apt-cache.
type noAptCommand struct{}

// NewNoAptCommand constructs the rule.
func NewNoAptCommand() engine.Rule { return noAptCommand{} }

// ID returns the rule identifier.
func (noAptCommand) ID() string { return "DL3027" }

// Check scans RUN instructions for apt usage.
func (noAptCommand) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		cmds := extractCommands(n)
		for _, c := range cmds {
			if c == "apt" {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3027",
					Message: "Do not use apt as it is meant to be an end-user tool, use apt-get or apt-cache instead",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}
