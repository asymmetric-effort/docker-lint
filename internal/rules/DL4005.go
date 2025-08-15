// file: internal/rules/DL4005.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// useShellForDefault warns when RUN is used to change the default shell.
type useShellForDefault struct{}

// NewUseShellForDefault constructs the rule.
func NewUseShellForDefault() engine.Rule { return useShellForDefault{} }

// ID returns the rule identifier.
func (useShellForDefault) ID() string { return "DL4005" }

// Check scans RUN instructions for ln commands targeting /bin/sh.
func (useShellForDefault) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		tokens := runTokens(n)
		cmds := splitTokens(tokens)
		for _, cmd := range cmds {
			if lnTargetsBinSh(cmd) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL4005",
					Message: "Use SHELL to change the default shell",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// lnTargetsBinSh reports ln invocations altering /bin/sh.
func lnTargetsBinSh(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	if strings.ToLower(tokens[0]) != "ln" {
		return false
	}
	for _, t := range tokens[1:] {
		if strings.Trim(t, "\"'") == "/bin/sh" {
			return true
		}
	}
	return false
}
