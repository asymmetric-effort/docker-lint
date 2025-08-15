package rules

/*
 * file: internal/rules/DL3056.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

var semverPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z-.]+)?(?:\+[0-9A-Za-z-.]+)?$`)

// labelSemVerValid ensures SemVer-typed labels conform to semantic versioning.
type labelSemVerValid struct{ schema LabelSchema }

// NewLabelSemVerValid constructs the rule.
func NewLabelSemVerValid(schema LabelSchema) engine.Rule { return &labelSemVerValid{schema: schema} }

// ID returns the rule identifier.
func (labelSemVerValid) ID() string { return "DL3056" }

// Check validates semantic version labels.
func (r *labelSemVerValid) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if r.schema[p.Key] == LabelTypeSemVer {
				if !semverPattern.MatchString(p.Value) {
					findings = append(findings, engine.Finding{RuleID: "DL3056", Message: "Label `" + p.Key + "` does not conform to semantic versioning.", Line: n.StartLine})
				}
			}
		}
	}
	return findings, nil
}
