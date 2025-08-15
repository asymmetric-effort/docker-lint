package rules

/*
 * file: internal/rules/DL3003.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// useWorkdir enforces using WORKDIR instead of cd in RUN.
type useWorkdir struct{}

// NewUseWorkdir constructs the rule.
func NewUseWorkdir() engine.Rule { return useWorkdir{} }

// ID returns the rule identifier.
func (useWorkdir) ID() string { return "DL3003" }

// Check scans RUN instructions for cd usage.
func (useWorkdir) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if c == "cd" {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3003",
					Message: "Use WORKDIR to switch to a directory",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}
