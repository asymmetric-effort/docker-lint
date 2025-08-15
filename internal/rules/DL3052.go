package rules

/*
 * file: internal/rules/DL3052.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"net/url"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// labelURLValid ensures URL-typed labels contain valid URLs.
type labelURLValid struct{ schema LabelSchema }

// NewLabelURLValid constructs the rule.
func NewLabelURLValid(schema LabelSchema) engine.Rule { return &labelURLValid{schema: schema} }

// ID returns the rule identifier.
func (labelURLValid) ID() string { return "DL3052" }

// Check validates URL labels against RFC 3986.
func (r *labelURLValid) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if r.schema[p.Key] == LabelTypeURL {
				u, err := url.Parse(p.Value)
				if err != nil || u.Scheme == "" || u.Host == "" {
					findings = append(findings, engine.Finding{RuleID: "DL3052", Message: "Label `" + p.Key + "` is not a valid URL.", Line: n.StartLine})
				}
			}
		}
	}
	return findings, nil
}
