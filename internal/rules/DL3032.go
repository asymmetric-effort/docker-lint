// file: internal/rules/DL3032.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireYumClean enforces cache cleanup after yum use.
type requireYumClean struct{}

// NewRequireYumClean constructs the rule.
func NewRequireYumClean() engine.Rule { return requireYumClean{} }

// ID returns the rule identifier.
func (requireYumClean) ID() string { return "DL3032" }

// Check ensures yum installs are followed by cleanup.
func (requireYumClean) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := splitRunSegments(n)
		if yumCleanMissing(segments) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3032",
				Message: "`yum clean all` missing after yum command.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// yumCleanMissing reports whether a yum install occurred without cleanup.
func yumCleanMissing(segs [][]string) bool {
	install := false
	clean := false
	for _, seg := range segs {
		if isYumInstall(seg) {
			install = true
		}
		if isYumClean(seg) {
			clean = true
		}
	}
	return install && !clean
}

// isYumClean detects yum clean or cache removal.
func isYumClean(seg []string) bool {
	if len(seg) >= 3 && seg[0] == "yum" && seg[1] == "clean" && seg[2] == "all" {
		return true
	}
	if len(seg) >= 3 && seg[0] == "rm" && seg[1] == "-rf" && seg[2] == "/var/cache/yum/*" {
		return true
	}
	return false
}
