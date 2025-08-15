// file: internal/rules/DL4004.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// singleEntrypoint ensures each build stage defines at most one ENTRYPOINT.
type singleEntrypoint struct{}

// NewSingleEntrypoint constructs the rule.
func NewSingleEntrypoint() engine.Rule { return singleEntrypoint{} }

// ID returns the rule identifier.
func (singleEntrypoint) ID() string { return "DL4004" }

// Check scans stages for multiple ENTRYPOINT instructions.
func (singleEntrypoint) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	hasEntrypoint := false
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "from") {
			hasEntrypoint = false
			continue
		}
		if strings.EqualFold(n.Value, "entrypoint") {
			if hasEntrypoint {
				findings = append(findings, engine.Finding{
					RuleID:  "DL4004",
					Message: "Multiple `ENTRYPOINT` instructions found. If you list more than one `ENTRYPOINT` then only the last `ENTRYPOINT` will take effect",
					Line:    n.StartLine,
				})
			} else {
				hasEntrypoint = true
			}
		}
	}
	return findings, nil
}
