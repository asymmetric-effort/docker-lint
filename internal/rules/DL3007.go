package rules

/*
 * file: internal/rules/DL3007.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"strings"
)

// noLatestTag ensures FROM instructions avoid implicit or latest tags.
type noLatestTag struct{}

// NewNoLatestTag constructs the rule.
func NewNoLatestTag() engine.Rule { return noLatestTag{} }

// ID returns the rule identifier.
func (noLatestTag) ID() string { return "DL3007" }

// Check evaluates each stage for usage of latest tags.
func (noLatestTag) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	for _, s := range d.Stages {
		if isLatest(s.From) {
			line := 0
			if s.Node != nil {
				line = s.Node.StartLine
			}
			findings = append(findings, engine.Finding{
				RuleID:  "DL3007",
				Message: "Using latest is prone to errors. Pin the version explicitly.",
				Line:    line,
			})
		}
	}
	return findings, nil
}

func isLatest(image string) bool {
	if strings.Contains(image, "@") {
		return false
	}
	parts := strings.Split(image, ":")
	if len(parts) == 1 {
		return true
	}
	return parts[len(parts)-1] == "latest"
}
