package rules

/*
 * file: internal/rules/DL3044.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// dnfVersionPin enforces version pinning for dnf/microdnf installs.
type dnfVersionPin struct{}

// NewDnfVersionPin constructs the rule.
func NewDnfVersionPin() engine.Rule { return dnfVersionPin{} }

// ID returns the rule identifier.
func (dnfVersionPin) ID() string { return "DL3044" }

// Check scans RUN instructions for unpinned dnf or microdnf installs.
func (dnfVersionPin) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := lowerSegments(splitRunSegments(n))
		for _, seg := range segments {
			if hasUnpinnedDnfInstall(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3044",
					Message: "Specify version with dnf/microdnf install. Use 'pkg-<version>' format for every installed package.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// hasUnpinnedDnfInstall reports whether tokens invoke dnf or microdnf install with unpinned packages.
func hasUnpinnedDnfInstall(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokens[0] != "dnf" && tokens[0] != "microdnf" {
		return false
	}
	idx := -1
	for i := 1; i < len(tokens); i++ {
		t := tokens[i]
		if t == "install" {
			idx = i
			break
		}
		if !strings.HasPrefix(t, "-") {
			return false
		}
	}
	if idx == -1 {
		return false
	}
	return unpinnedDnfPackages(tokens[idx+1:])
}

// unpinnedDnfPackages checks package arguments for version pins.
func unpinnedDnfPackages(args []string) bool {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		dash := strings.Index(a, "-")
		if dash <= 0 {
			return true
		}
		ver := a[dash+1:]
		if ver == "" || strings.HasPrefix(ver, "-") {
			return true
		}
	}
	return false
}
