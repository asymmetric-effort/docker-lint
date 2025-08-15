// file: internal/rules/DL3024.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// uniqueStageNames ensures FROM aliases are unique.
type uniqueStageNames struct{}

// NewUniqueStageNames constructs the rule.
func NewUniqueStageNames() engine.Rule { return uniqueStageNames{} }

// ID returns the rule identifier.
func (uniqueStageNames) ID() string { return "DL3024" }

// Check verifies no duplicate stage aliases exist.
func (uniqueStageNames) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	aliases := map[string]int{}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "from") {
			continue
		}
		if name := strings.ToLower(stageAlias(n)); name != "" {
			if line, ok := aliases[name]; ok {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3024",
					Message: "FROM aliases (stage names) must be unique",
					Line:    n.StartLine,
				})
				_ = line // previously stored; not used further
			} else {
				aliases[name] = n.StartLine
			}
		}
	}
	return findings, nil
}
