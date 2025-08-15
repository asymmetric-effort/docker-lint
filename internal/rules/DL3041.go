package rules

/*
 * file: internal/rules/DL3041.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// dnfNoUpgrade warns against using dnf or microdnf upgrade/update.
type dnfNoUpgrade struct{}

// NewDnfNoUpgrade constructs the rule.
func NewDnfNoUpgrade() engine.Rule { return dnfNoUpgrade{} }

// ID returns the rule identifier.
func (dnfNoUpgrade) ID() string { return "DL3041" }

// Check scans RUN instructions for disallowed dnf upgrade or update usage.
func (dnfNoUpgrade) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segs := lowerSegments(splitRunSegments(n))
		for _, seg := range segs {
			if isDnfUpgrade(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3041",
					Message: "Avoid dnf upgrade or update; use a newer base image or install specific packages with pinned versions instead.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isDnfUpgrade reports whether the segment invokes dnf or microdnf upgrade/update.
func isDnfUpgrade(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokens[0] != "dnf" && tokens[0] != "microdnf" {
		return false
	}
	for _, t := range tokens[1:] {
		if strings.HasPrefix(t, "-") {
			continue
		}
		return t == "upgrade" || t == "update"
	}
	return false
}
