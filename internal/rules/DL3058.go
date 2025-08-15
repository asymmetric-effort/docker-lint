package rules

/*
 * file: internal/rules/DL3058.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"net/mail"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// labelEmailValid ensures email-typed labels contain valid addresses.
type labelEmailValid struct{ schema LabelSchema }

// NewLabelEmailValid constructs the rule.
func NewLabelEmailValid(schema LabelSchema) engine.Rule { return &labelEmailValid{schema: schema} }

// ID returns the rule identifier.
func (labelEmailValid) ID() string { return "DL3058" }

// Check validates email label values.
func (r *labelEmailValid) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if r.schema[p.Key] == LabelTypeEmail {
				if _, err := mail.ParseAddress(p.Value); err != nil {
					findings = append(findings, engine.Finding{RuleID: "DL3058", Message: "Label `" + p.Key + "` is not a valid email format - must conform to RFC5322.", Line: n.StartLine})
				}
			}
		}
	}
	return findings, nil
}
