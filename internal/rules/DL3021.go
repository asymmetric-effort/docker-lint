package rules

/*
 * file: internal/rules/DL3021.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// copyDestEndsWithSlash ensures COPY with multiple sources uses a directory destination.
type copyDestEndsWithSlash struct{}

// NewCopyDestEndsWithSlash constructs the rule.
func NewCopyDestEndsWithSlash() engine.Rule { return copyDestEndsWithSlash{} }

// ID returns the rule identifier.
func (copyDestEndsWithSlash) ID() string { return "DL3021" }

// Check verifies that COPY with more than 2 arguments uses a destination ending with '/'.
func (copyDestEndsWithSlash) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "copy") {
			continue
		}
		tokens := collectArgs(n)
		if len(tokens) <= 2 {
			continue
		}
		dest := tokens[len(tokens)-1]
		if !strings.HasSuffix(dest, "/") {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3021",
				Message: "COPY with more than 2 arguments requires the last argument to end with /",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}
