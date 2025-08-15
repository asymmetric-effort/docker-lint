package rules

/*
 * file: internal/rules/DL1001.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// noInlineIgnore detects inline hadolint ignore pragmas.
type noInlineIgnore struct{}

// NewNoInlineIgnore constructs the rule.
func NewNoInlineIgnore() engine.Rule { return noInlineIgnore{} }

// ID returns the rule identifier.
func (noInlineIgnore) ID() string { return "DL1001" }

// Check scans comments and instructions for hadolint ignore pragmas.
func (noInlineIgnore) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		for _, com := range n.PrevComment {
			if hasIgnorePragma(com) {
				line := n.StartLine
				if line > 0 {
					line -= len(n.PrevComment)
				}
				findings = append(findings, engine.Finding{
					RuleID:  "DL1001",
					Message: "Please refrain from using inline ignore pragmas `# hadolint ignore=DLxxxx`.",
					Line:    line,
				})
			}
		}
		if hasIgnorePragma(n.Original) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL1001",
				Message: "Please refrain from using inline ignore pragmas `# hadolint ignore=DLxxxx`.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// hasIgnorePragma reports if the string contains a hadolint ignore pragma.
func hasIgnorePragma(s string) bool {
	return strings.Contains(strings.ToLower(s), "hadolint ignore=")
}
