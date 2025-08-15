package rules

/*
 * file: internal/rules/DL3042.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// combinePackageRuns flags consecutive RUN instructions using the same package manager.
type combinePackageRuns struct{}

// NewCombinePackageRuns constructs the rule.
func NewCombinePackageRuns() engine.Rule { return combinePackageRuns{} }

// ID returns the rule identifier.
func (combinePackageRuns) ID() string { return "DL3042" }

// Check detects consecutive package manager RUN instructions.
func (combinePackageRuns) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	prev := ""
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			if strings.EqualFold(n.Value, "#") {
				continue
			}
			prev = ""
			continue
		}
		pm := runPackageManager(n)
		if pm != "" && pm == prev {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3042",
				Message: "Combine consecutive RUN instructions that use the same package manager.",
				Line:    n.StartLine,
			})
		}
		prev = pm
	}
	return findings, nil
}

// runPackageManager determines the package manager family used in a RUN instruction.
func runPackageManager(n *parser.Node) string {
	fams := make(map[string]struct{})
	for _, seg := range splitRunSegments(n) {
		if fam := packageManagerFamily(seg); fam != "" {
			fams[fam] = struct{}{}
		}
	}
	if len(fams) != 1 {
		return ""
	}
	for k := range fams {
		return k
	}
	return ""
}

// packageManagerFamily returns the package manager family for a command segment.
func packageManagerFamily(seg []string) string {
	if len(seg) < 2 {
		return ""
	}
	cmd := seg[0]
	i := 1
	for i < len(seg) && strings.HasPrefix(seg[i], "-") {
		i++
	}
	if i >= len(seg) {
		return ""
	}
	sub := seg[i]
	switch cmd {
	case "apt-get", "apt":
		switch sub {
		case "update", "install", "upgrade", "remove", "clean":
			return "apt"
		}
	case "apk":
		switch sub {
		case "add", "del", "update", "upgrade", "fix", "cache":
			return "apk"
		}
	case "dnf", "microdnf":
		switch sub {
		case "update", "install", "upgrade", "remove", "clean":
			return "dnf"
		}
	case "yum":
		switch sub {
		case "update", "install", "upgrade", "remove", "clean":
			return "yum"
		}
	case "zypper":
		switch sub {
		case "update", "install", "upgrade", "remove", "clean":
			return "zypper"
		}
	}
	return ""
}
