package rules

/*
 * file: internal/rules/DL3054.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

var spdxPattern = regexp.MustCompile(`^[A-Za-z0-9-.+]+$`)

// labelSPDXValid ensures SPDX-typed labels match SPDX identifier pattern.
type labelSPDXValid struct{ schema LabelSchema }

// NewLabelSPDXValid constructs the rule.
func NewLabelSPDXValid(schema LabelSchema) engine.Rule { return &labelSPDXValid{schema: schema} }

// ID returns the rule identifier.
func (labelSPDXValid) ID() string { return "DL3054" }

// Check validates SPDX label values.
func (r *labelSPDXValid) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if r.schema[p.Key] == LabelTypeSPDX {
				if !spdxPattern.MatchString(p.Value) {
					findings = append(findings, engine.Finding{RuleID: "DL3054", Message: "Label `" + p.Key + "` is not a valid SPDX identifier.", Line: n.StartLine})
				}
			}
		}
	}
	return findings, nil
}
