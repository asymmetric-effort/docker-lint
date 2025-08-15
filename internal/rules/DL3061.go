// file: internal/rules/DL3061.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// startWithFromOrArg ensures the Dockerfile begins with FROM or ARG.
type startWithFromOrArg struct{}

// NewStartWithFromOrArg constructs the rule.
func NewStartWithFromOrArg() engine.Rule { return startWithFromOrArg{} }

// ID returns the rule identifier.
func (startWithFromOrArg) ID() string { return "DL3061" }

// Check verifies that the first instruction is FROM or ARG.
func (startWithFromOrArg) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil || len(d.AST.Children) == 0 {
		return findings, nil
	}
	first := d.AST.Children[0]
	if !strings.EqualFold(first.Value, "from") && !strings.EqualFold(first.Value, "arg") {
		findings = append(findings, engine.Finding{
			RuleID:  "DL3061",
			Message: "Invalid instruction order. Dockerfile must begin with `FROM`, `ARG` or comment.",
			Line:    first.StartLine,
		})
	}
	return findings, nil
}
