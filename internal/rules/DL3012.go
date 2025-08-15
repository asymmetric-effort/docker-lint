package rules

/*
 * file: internal/rules/DL3012.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// singleHealthcheck ensures only one HEALTHCHECK instruction per stage.
type singleHealthcheck struct{}

// NewSingleHealthcheck constructs the rule.
func NewSingleHealthcheck() engine.Rule { return singleHealthcheck{} }

// ID returns the rule identifier.
func (singleHealthcheck) ID() string { return "DL3012" }

// Check verifies that each stage contains at most one HEALTHCHECK instruction.
func (singleHealthcheck) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
		if strings.EqualFold(n.Value, "healthcheck") {
			if seen {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3012",
					Message: "Multiple HEALTHCHECK instructions",
					Line:    n.StartLine,
				})
			} else {
				seen = true
			}
		}
	}
	return findings, nil
}
