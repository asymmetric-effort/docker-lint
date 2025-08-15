package rules

/*
 * file: internal/rules/DL3051.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// labelNotEmpty ensures specific labels are not empty.
type labelNotEmpty struct{ schema LabelSchema }

// NewLabelNotEmpty constructs the rule.
func NewLabelNotEmpty(schema LabelSchema) engine.Rule { return &labelNotEmpty{schema: schema} }

// ID returns the rule identifier.
func (labelNotEmpty) ID() string { return "DL3051" }

// Check reports schema-defined labels that have empty values.
func (r *labelNotEmpty) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if inSchema(r.schema, p.Key) && p.Value == "" {
				findings = append(findings, engine.Finding{RuleID: "DL3051", Message: "label `" + p.Key + "` is empty.", Line: n.StartLine})
			}
		}
	}
	return findings, nil
}
