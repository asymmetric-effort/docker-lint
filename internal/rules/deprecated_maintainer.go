// file: internal/rules/deprecated_maintainer.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// deprecatedMaintainer flags usage of the deprecated MAINTAINER instruction.
type deprecatedMaintainer struct{}

// NewDeprecatedMaintainer constructs the rule.
func NewDeprecatedMaintainer() engine.Rule { return deprecatedMaintainer{} }

// ID returns the rule identifier.
func (deprecatedMaintainer) ID() string { return "DL4000" }

// Check scans the AST for MAINTAINER instructions.
func (deprecatedMaintainer) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "maintainer") {
			line := n.StartLine
			findings = append(findings, engine.Finding{
				RuleID:  "DL4000",
				Message: "MAINTAINER is deprecated. Use LABEL maintainer=\"name\" instead.",
				Line:    line,
			})
		}
	}
	return findings, nil
}
