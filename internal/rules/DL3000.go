package rules

/*
 * file: internal/rules/DL3000.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"
	"unicode"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// absoluteWorkdir ensures WORKDIR uses absolute paths.
type absoluteWorkdir struct{}

// NewAbsoluteWorkdir constructs the rule.
func NewAbsoluteWorkdir() engine.Rule { return absoluteWorkdir{} }

// ID returns the rule identifier.
func (absoluteWorkdir) ID() string { return "DL3000" }

// Check examines WORKDIR instructions for absolute paths.
func (absoluteWorkdir) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "workdir") {
			line := n.StartLine
			path := ""
			if n.Next != nil {
				path = n.Next.Value
			}
			p := strings.Trim(path, "\"'")
			if !isAbsoluteWorkdir(p) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3000",
					Message: "Use absolute WORKDIR",
					Line:    line,
				})
			}
		}
	}
	return findings, nil
}

// isAbsoluteWorkdir reports whether the provided path is absolute or variable-based.
func isAbsoluteWorkdir(p string) bool {
	if strings.HasPrefix(p, "$") {
		return true
	}
	if strings.HasPrefix(p, "/") {
		return true
	}
	if len(p) > 1 && unicode.IsLetter(rune(p[0])) && p[1] == ':' {
		return true
	}
	return false
}
