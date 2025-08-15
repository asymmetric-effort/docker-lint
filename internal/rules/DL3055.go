package rules

/*
 * file: internal/rules/DL3055.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

var gitHashPattern = regexp.MustCompile(`^[0-9a-f]{7}([0-9a-f]{33})?$`)

// labelGitHashValid ensures Git hash labels are valid.
type labelGitHashValid struct{ schema LabelSchema }

// NewLabelGitHashValid constructs the rule.
func NewLabelGitHashValid(schema LabelSchema) engine.Rule { return &labelGitHashValid{schema: schema} }

// ID returns the rule identifier.
func (labelGitHashValid) ID() string { return "DL3055" }

// Check validates Git hash label values.
func (r *labelGitHashValid) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "label") {
			continue
		}
		for _, p := range collectLabelPairs(n) {
			if r.schema[p.Key] == LabelTypeGitHash {
				if !gitHashPattern.MatchString(strings.ToLower(p.Value)) {
					findings = append(findings, engine.Finding{RuleID: "DL3055", Message: "Label `" + p.Key + "` is not a valid git hash.", Line: n.StartLine})
				}
			}
		}
	}
	return findings, nil
}
