package rules

/*
 * file: internal/rules/DL3050.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// superfluousLabels checks for labels not defined in the schema when strict mode is enabled.
type superfluousLabels struct {
	schema LabelSchema
	strict bool
}

// NewSuperfluousLabels constructs the rule.
func NewSuperfluousLabels(schema LabelSchema, strict bool) engine.Rule {
	return &superfluousLabels{schema: schema, strict: strict}
}

// ID returns the rule identifier.
func (superfluousLabels) ID() string { return "DL3050" }

// Check reports labels not present in the schema when strict mode is enabled.
func (r *superfluousLabels) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if !r.strict || d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if !inSchema(r.schema, p.Key) {
				findings = append(findings, engine.Finding{RuleID: "DL3050", Message: "Superfluous label(s) present.", Line: n.StartLine})
				break
			}
		}
	}
	return findings, nil
}
