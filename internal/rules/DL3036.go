// file: internal/rules/DL3036.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireZypperClean enforces cache cleanup after zypper use.
type requireZypperClean struct{}

// NewRequireZypperClean constructs the rule.
func NewRequireZypperClean() engine.Rule { return requireZypperClean{} }

// ID returns the rule identifier.
func (requireZypperClean) ID() string { return "DL3036" }

// Check ensures zypper installs are followed by cleanup.
func (requireZypperClean) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := splitRunSegments(n)
		if zypperCleanMissing(segments) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3036",
				Message: "`zypper clean` missing after zypper use.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// zypperCleanMissing reports if install without cleanup occurs.
func zypperCleanMissing(segs [][]string) bool {
	install := false
	clean := false
	for _, seg := range segs {
		if isZypperInstall(seg) {
			install = true
		}
		if isZypperClean(seg) {
			clean = true
		}
	}
	return install && !clean
}

// isZypperInstall detects zypper install/in.
func isZypperInstall(seg []string) bool {
	if len(seg) < 2 || seg[0] != "zypper" {
		return false
	}
	if seg[1] == "install" || seg[1] == "in" {
		return true
	}
	return false
}

// isZypperClean detects zypper clean/cc.
func isZypperClean(seg []string) bool {
	if len(seg) >= 2 && seg[0] == "zypper" && (seg[1] == "clean" || seg[1] == "cc") {
		return true
	}
	return false
}
