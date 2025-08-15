// file: internal/rules/DL3037.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// pinZypperVersions enforces version pinning for zypper installs.
type pinZypperVersions struct{}

// NewPinZypperVersions constructs the rule.
func NewPinZypperVersions() engine.Rule { return pinZypperVersions{} }

// ID returns the rule identifier.
func (pinZypperVersions) ID() string { return "DL3037" }

// Check scans RUN instructions for unpinned zypper packages.
func (pinZypperVersions) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if isZypperInstall(seg) {
				pkgs := collectNonFlag(seg[2:])
				if len(pkgs) > 0 && !allZypperVersionFixed(pkgs) {
					findings = append(findings, engine.Finding{
						RuleID:  "DL3037",
						Message: "Specify version with `zypper install -y <package>=<version>`.",
						Line:    n.StartLine,
					})
					break
				}
			}
		}
	}
	return findings, nil
}

// allZypperVersionFixed returns true if each package pins a version.
func allZypperVersionFixed(pkgs []string) bool {
	for _, p := range pkgs {
		if !(strings.Contains(p, "=") || strings.Contains(p, ">=") || strings.Contains(p, ">") || strings.Contains(p, "<=") || strings.Contains(p, "<") || strings.HasSuffix(p, ".rpm")) {
			return false
		}
	}
	return true
}
