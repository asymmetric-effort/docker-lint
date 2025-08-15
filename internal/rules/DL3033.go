// file: internal/rules/DL3033.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// pinYumVersions enforces version pinning in yum installs.
type pinYumVersions struct{}

// NewPinYumVersions constructs the rule.
func NewPinYumVersions() engine.Rule { return pinYumVersions{} }

// ID returns the rule identifier.
func (pinYumVersions) ID() string { return "DL3033" }

// Check scans RUN instructions for unpinned yum packages.
func (pinYumVersions) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
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
			if len(seg) < 2 || seg[0] != "yum" {
				continue
			}
			if violatesYumPin(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3033",
					Message: "Specify version with `yum install -y <package>-<version>`.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// violatesYumPin reports if yum packages or modules lack versions.
func violatesYumPin(seg []string) bool {
	if seg[1] == "module" {
		idx := indexOf(seg, "install")
		if idx == -1 {
			return false
		}
		pkgs := collectNonFlag(seg[idx+1:])
		for _, p := range pkgs {
			if !strings.Contains(p, ":") {
				return true
			}
		}
		return false
	}
	// regular install
	idx := indexOf(seg, "install")
	if idx == -1 {
		return false
	}
	pkgs := collectNonFlag(seg[idx+1:])
	for _, p := range pkgs {
		if !strings.Contains(p, "-") && !strings.HasSuffix(p, ".rpm") {
			return true
		}
	}
	return false
}

// indexOf returns index of token or -1.
func indexOf(slice []string, target string) int {
	for i, t := range slice {
		if t == target {
			return i
		}
	}
	return -1
}

// collectNonFlag gathers non-flag arguments from token list.
func collectNonFlag(tokens []string) []string {
	var out []string
	for _, t := range tokens {
		if strings.HasPrefix(t, "-") {
			continue
		}
		out = append(out, t)
	}
	return out
}
