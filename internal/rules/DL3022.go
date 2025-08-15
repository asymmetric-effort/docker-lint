// file: internal/rules/DL3022.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// copyFromPreviousStage verifies COPY --from references a prior stage.
type copyFromPreviousStage struct{}

// NewCopyFromPreviousStage constructs the rule.
func NewCopyFromPreviousStage() engine.Rule { return copyFromPreviousStage{} }

// ID returns the rule identifier.
func (copyFromPreviousStage) ID() string { return "DL3022" }

// Check ensures COPY --from references a previously defined stage alias or index.
func (copyFromPreviousStage) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	aliases := map[string]struct{}{}
	stageCount := 0
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "from") {
			stageCount++
			if name := stageAlias(n); name != "" {
				aliases[strings.ToLower(name)] = struct{}{}
			}
			continue
		}
		if !strings.EqualFold(n.Value, "copy") {
			continue
		}
		from, ok := copyFromFlag(n)
		if !ok {
			continue
		}
		if strings.Contains(from, ":") {
			continue
		}
		if _, ok := aliases[strings.ToLower(from)]; ok {
			continue
		}
		if idx, err := strconv.Atoi(from); err == nil {
			if idx < stageCount-1 {
				continue
			}
		}
		findings = append(findings, engine.Finding{
			RuleID:  "DL3022",
			Message: "`COPY --from` should reference a previously defined `FROM` alias",
			Line:    n.StartLine,
		})
	}
	return findings, nil
}
