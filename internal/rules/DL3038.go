// file: internal/rules/DL3038.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireDnfYes enforces non-interactive dnf commands.
type requireDnfYes struct{}

// NewRequireDnfYes constructs the rule.
func NewRequireDnfYes() engine.Rule { return requireDnfYes{} }

// ID returns the rule identifier.
func (requireDnfYes) ID() string { return "DL3038" }

// Check inspects RUN instructions for dnf installs without -y.
func (requireDnfYes) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if isDnfInstall(seg) && !hasDnfYes(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3038",
					Message: "Use the -y switch to avoid manual input `dnf install -y <package>`",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isDnfInstall reports whether segment invokes dnf/microdnf install.
func isDnfInstall(seg []string) bool {
	if len(seg) < 2 {
		return false
	}
	if seg[0] != "dnf" && seg[0] != "microdnf" {
		return false
	}
	for _, t := range seg[1:] {
		if t == "install" || t == "groupinstall" || t == "localinstall" {
			return true
		}
	}
	return false
}

// hasDnfYes detects non-interactive flags.
func hasDnfYes(seg []string) bool {
	for _, t := range seg {
		if t == "-y" || t == "--assumeyes" || strings.HasPrefix(t, "--assumeyes=") {
			return true
		}
	}
	return false
}
