package rules

/*
 * file: internal/rules/DL3046.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// apkNoUpgrade discourages using `apk upgrade` in Dockerfiles.
type apkNoUpgrade struct{}

// NewApkNoUpgrade constructs the rule.
func NewApkNoUpgrade() engine.Rule { return apkNoUpgrade{} }

// ID returns the rule identifier.
func (apkNoUpgrade) ID() string { return "DL3046" }

// Check scans RUN instructions for `apk upgrade` invocations.
func (apkNoUpgrade) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segs := lowerSegments(splitRunSegments(n))
		for _, seg := range segs {
			if isApkUpgrade(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3046",
					Message: "Avoid apk upgrade in Dockerfiles. Upgrade the base image or install specific pinned packages instead.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isApkUpgrade reports whether the segment invokes `apk upgrade`.
func isApkUpgrade(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}
	if tokens[0] != "apk" {
		return false
	}
	for _, t := range tokens[1:] {
		if strings.HasPrefix(t, "-") {
			continue
		}
		return t == "upgrade"
	}
	return false
}
