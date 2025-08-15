package rules

/*
 * file: internal/rules/DL3011.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// validPortRange ensures exposed ports are within valid range.
type validPortRange struct{}

// NewValidPortRange constructs the rule.
func NewValidPortRange() engine.Rule { return validPortRange{} }

// ID returns the rule identifier.
func (validPortRange) ID() string { return "DL3011" }

// Check inspects EXPOSE instructions for invalid ports.
func (validPortRange) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "expose") {
			continue
		}
		for arg := n.Next; arg != nil; arg = arg.Next {
			if !portInRange(arg.Value) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3011",
					Message: "Valid UNIX ports range from 0 to 65535",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// portInRange reports whether all numeric parts of a port token are within 0-65535.
func portInRange(token string) bool {
	token = strings.SplitN(token, "/", 2)[0]
	parts := strings.Split(token, "-")
	for _, p := range parts {
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return true
		}
		if v < 0 || v > 65535 {
			return false
		}
	}
	return true
}
