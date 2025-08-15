// file: internal/rules/DL3029.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// noPlatformInFrom disallows explicit --platform in FROM.
type noPlatformInFrom struct{}

// NewNoPlatformInFrom constructs the rule.
func NewNoPlatformInFrom() engine.Rule { return noPlatformInFrom{} }

// ID returns the rule identifier.
func (noPlatformInFrom) ID() string { return "DL3029" }

// Check warns on --platform usage in FROM instructions.
func (noPlatformInFrom) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "from") {
			continue
		}
		for _, f := range n.Flags {
			if strings.HasPrefix(strings.ToLower(f), "--platform=") {
				v := strings.TrimPrefix(f, "--platform=")
				v = strings.Trim(v, "\"'")
				if !strings.Contains(v, "BUILDPLATFORM") && !strings.Contains(v, "TARGETPLATFORM") {
					findings = append(findings, engine.Finding{
						RuleID:  "DL3029",
						Message: "Do not use --platform flag with FROM",
						Line:    n.StartLine,
					})
				}
				break
			}
		}
	}
	return findings, nil
}
