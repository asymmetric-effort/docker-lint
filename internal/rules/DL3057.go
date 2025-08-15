package rules

/*
 * file: internal/rules/DL3057.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// healthcheckExists reports stages missing a HEALTHCHECK instruction.
type healthcheckExists struct{}

// NewHealthcheckExists constructs the rule.
func NewHealthcheckExists() engine.Rule { return healthcheckExists{} }

// ID returns the rule identifier.
func (healthcheckExists) ID() string { return "DL3057" }

// Check verifies that each stage or its base image defines a HEALTHCHECK.
func (healthcheckExists) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	type stage struct {
		line     int
		parent   int
		hasCheck bool
	}
	var stages []*stage
	nameMap := make(map[string]int)
	current := -1
	for _, n := range d.AST.Children {
		switch strings.ToLower(n.Value) {
		case "from":
			base := ""
			alias := ""
			if n.Next != nil {
				base = n.Next.Value
				for tok := n.Next.Next; tok != nil; tok = tok.Next {
					if strings.EqualFold(tok.Value, "as") && tok.Next != nil {
						alias = tok.Next.Value
						break
					}
				}
			}
			parent := -1
			if idx, ok := nameMap[base]; ok {
				parent = idx
			}
			st := &stage{line: n.StartLine, parent: parent}
			stages = append(stages, st)
			name := alias
			if name == "" {
				name = base
			}
			nameMap[name] = len(stages) - 1
			current = len(stages) - 1
		case "healthcheck":
			if current >= 0 {
				stages[current].hasCheck = true
			}
		}
	}
	cache := make(map[int]bool)
	var hasHC func(int) bool
	hasHC = func(id int) bool {
		if v, ok := cache[id]; ok {
			return v
		}
		st := stages[id]
		if st.hasCheck {
			cache[id] = true
			return true
		}
		if st.parent == -1 {
			cache[id] = false
			return false
		}
		v := hasHC(st.parent)
		cache[id] = v
		return v
	}
	for id, st := range stages {
		if !hasHC(id) {
			findings = append(findings, engine.Finding{RuleID: "DL3057", Message: "`HEALTHCHECK` instruction missing.", Line: st.line})
		}
	}
	return findings, nil
}
