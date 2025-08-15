package rules

/*
 * file: internal/rules/DL3053.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"
	"time"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// labelTimeRFC3339 ensures RFC3339-typed labels are valid timestamps.
type labelTimeRFC3339 struct{ schema LabelSchema }

// NewLabelTimeRFC3339 constructs the rule.
func NewLabelTimeRFC3339(schema LabelSchema) engine.Rule { return &labelTimeRFC3339{schema: schema} }

// ID returns the rule identifier.
func (labelTimeRFC3339) ID() string { return "DL3053" }

// Check validates time labels.
func (r *labelTimeRFC3339) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if r.schema[p.Key] == LabelTypeRFC3339 {
				if _, err := time.Parse(time.RFC3339, p.Value); err != nil {
					findings = append(findings, engine.Finding{RuleID: "DL3053", Message: "Label `" + p.Key + "` is not a valid time format - must conform to RFC3339.", Line: n.StartLine})
				}
			}
		}
	}
	return findings, nil
}
