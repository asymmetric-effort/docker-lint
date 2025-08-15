package rules

/*
 * file: internal/rules/DL3006.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireTag enforces explicit image tags in FROM instructions.
type requireTag struct{}

// NewRequireTag constructs the rule.
func NewRequireTag() engine.Rule { return requireTag{} }

// ID returns the rule identifier.
func (requireTag) ID() string { return "DL3006" }

// Check examines stage base images for explicit tags.
func (requireTag) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil {
		return findings, nil
	}
	aliases := map[string]struct{}{}
	for _, s := range d.Stages {
		from := s.From
		if from == "" {
			if s.Name != "" {
				aliases[s.Name] = struct{}{}
			}
			continue
		}
		_, isAlias := aliases[from]
		if needsTag(from, isAlias) {
			line := 0
			if s.Node != nil {
				line = s.Node.StartLine
			}
			findings = append(findings, engine.Finding{
				RuleID:  "DL3006",
				Message: "Always tag the version of an image explicitly",
				Line:    line,
			})
		}
		if s.Name != "" {
			aliases[s.Name] = struct{}{}
		}
	}
	return findings, nil
}

// needsTag determines whether the FROM image is missing a tag or digest.
func needsTag(image string, isAlias bool) bool {
	if isAlias {
		return false
	}
	if image == "scratch" {
		return false
	}
	if strings.HasPrefix(image, "$") {
		return false
	}
	if strings.Contains(image, "@") {
		return false
	}
	parts := strings.Split(image, ":")
	return len(parts) == 1
}
