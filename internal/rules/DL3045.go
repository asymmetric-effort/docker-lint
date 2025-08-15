// file: internal/rules/DL3045.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// copyFromExternalDigest ensures COPY --from external images include a digest.
type copyFromExternalDigest struct{}

// NewCopyFromExternalDigest constructs the rule.
func NewCopyFromExternalDigest() engine.Rule { return copyFromExternalDigest{} }

// ID returns the rule identifier.
func (copyFromExternalDigest) ID() string { return "DL3045" }

// Check flags external COPY --from references lacking a digest.
func (copyFromExternalDigest) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	aliases := map[string]struct{}{}
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "from") {
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
		lf := strings.ToLower(from)
		if lf == "scratch" {
			continue
		}
		if strings.HasPrefix(from, "$") {
			continue
		}
		if _, ok := aliases[lf]; ok {
			continue
		}
		if _, err := strconv.Atoi(from); err == nil {
			continue
		}
		if !strings.Contains(lf, "@sha256:") {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3045",
				Message: "COPY --from without digest pinning for external image.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}
