package rules

/*
 * file: internal/rules/DL3002.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// lastUserNotRoot ensures the final USER in each stage is non-root.
type lastUserNotRoot struct{}

// NewLastUserNotRoot constructs the rule.
func NewLastUserNotRoot() engine.Rule { return lastUserNotRoot{} }

// ID returns the rule identifier.
func (lastUserNotRoot) ID() string { return "DL3002" }

// Check verifies that the last USER instruction per stage is non-root.
func (lastUserNotRoot) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	children := d.AST.Children
	for i, stage := range d.Stages {
		start := -1
		end := len(children)
		for idx, n := range children {
			if n == stage.Node {
				start = idx
				break
			}
		}
		if start == -1 {
			continue
		}
		if i+1 < len(d.Stages) {
			for idx, n := range children {
				if n == d.Stages[i+1].Node {
					end = idx
					break
				}
			}
		}
		last := ""
		line := 0
		for _, n := range children[start+1 : end] {
			if strings.EqualFold(n.Value, "user") && n.Next != nil {
				last = n.Next.Value
				line = n.StartLine
			}
		}
		if last != "" && isRootUser(last) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3002",
				Message: "Last USER should not be root",
				Line:    line,
			})
		}
	}
	return findings, nil
}

func isRootUser(u string) bool {
	return u == "root" || u == "0" || strings.HasPrefix(u, "root:") || strings.HasPrefix(u, "0:")
}
