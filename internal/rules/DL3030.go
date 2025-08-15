// file: internal/rules/DL3030.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireYumYes enforces non-interactive yum commands.
type requireYumYes struct{}

// NewRequireYumYes constructs the rule.
func NewRequireYumYes() engine.Rule { return requireYumYes{} }

// ID returns the rule identifier.
func (requireYumYes) ID() string { return "DL3030" }

// Check inspects RUN instructions for yum install without -y.
func (requireYumYes) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := splitRunSegments(n)
		for _, seg := range segments {
			if isYumInstall(seg) && !hasYumYes(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3030",
					Message: "Use the -y switch to avoid manual input `yum install -y <package>`",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isYumInstall reports whether the segment invokes yum install variants.
func isYumInstall(seg []string) bool {
	if len(seg) < 2 || seg[0] != "yum" {
		return false
	}
	for _, t := range seg[1:] {
		if t == "install" || t == "groupinstall" || t == "localinstall" {
			return true
		}
	}
	return false
}

// hasYumYes reports if non-interactive flags are present.
func hasYumYes(seg []string) bool {
	for _, t := range seg {
		if t == "-y" || t == "--assumeyes" || strings.HasPrefix(t, "--assumeyes=") {
			return true
		}
	}
	return false
}
