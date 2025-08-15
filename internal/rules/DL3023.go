// file: internal/rules/DL3023.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// copyFromSelf disallows COPY --from referencing current stage.
type copyFromSelf struct{}

// NewCopyFromSelf constructs the rule.
func NewCopyFromSelf() engine.Rule { return copyFromSelf{} }

// ID returns the rule identifier.
func (copyFromSelf) ID() string { return "DL3023" }

// Check flags COPY --from that references its own FROM alias or index.
func (copyFromSelf) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	currentAlias := ""
	currentIndex := -1
	index := -1
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "from") {
			index++
			currentIndex = index
			currentAlias = strings.ToLower(stageAlias(n))
			continue
		}
		if !strings.EqualFold(n.Value, "copy") {
			continue
		}
		from, ok := copyFromFlag(n)
		if !ok {
			continue
		}
		if strings.EqualFold(strings.ToLower(from), currentAlias) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3023",
				Message: "`COPY --from` cannot reference its own `FROM` alias",
				Line:    n.StartLine,
			})
			continue
		}
		if idx, err := strconv.Atoi(from); err == nil {
			if idx == currentIndex {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3023",
					Message: "`COPY --from` cannot reference its own `FROM` alias",
					Line:    n.StartLine,
				})
			}
		}
	}
	return findings, nil
}
