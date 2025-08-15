package rules

/*
 * file: internal/rules/DL3004.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// noSudo disallows sudo in RUN instructions.
type noSudo struct{}

// NewNoSudo constructs the rule.
func NewNoSudo() engine.Rule { return noSudo{} }

// ID returns the rule identifier.
func (noSudo) ID() string { return "DL3004" }

// Check scans RUN instructions for sudo usage.
func (noSudo) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if c == "sudo" {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3004",
					Message: "Do not use sudo as it leads to unpredictable behavior. Use a tool like gosu to enforce root",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}
