// file: internal/rules/DL3025.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// jsonNotationCmdEntrypoint enforces JSON array form for CMD and ENTRYPOINT.
type jsonNotationCmdEntrypoint struct{}

// NewJSONNotationCmdEntrypoint constructs the rule.
func NewJSONNotationCmdEntrypoint() engine.Rule { return jsonNotationCmdEntrypoint{} }

// ID returns the rule identifier.
func (jsonNotationCmdEntrypoint) ID() string { return "DL3025" }

// Check reports CMD or ENTRYPOINT using shell form.
func (jsonNotationCmdEntrypoint) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "cmd") && !strings.EqualFold(n.Value, "entrypoint") {
			continue
		}
		if n.Attributes == nil || !n.Attributes["json"] {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3025",
				Message: "Use arguments JSON notation for CMD and ENTRYPOINT arguments",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}
