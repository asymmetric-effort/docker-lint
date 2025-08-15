package rules

/*
 * file: internal/rules/DL3040.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// dnfCacheCleanup ensures dnf or microdnf package operations clean metadata in the same layer.
type dnfCacheCleanup struct{}

// NewDnfCacheCleanup constructs the rule.
func NewDnfCacheCleanup() engine.Rule { return dnfCacheCleanup{} }

// ID returns the rule identifier.
func (dnfCacheCleanup) ID() string { return "DL3040" }

// Check scans RUN instructions for dnf/microdnf commands lacking cleanup.
func (dnfCacheCleanup) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := lowerSegments(splitRunSegments(n))
		if needsDnfCleanup(segments) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3040",
				Message: "dnf clean all missing after dnf command.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// needsDnfCleanup reports whether a dnf/microdnf operation lacks subsequent cleanup.
func needsDnfCleanup(segments [][]string) bool {
	lastOp := -1
	lastClean := -1
	for i, seg := range segments {
		if isDnfOperation(seg) {
			lastOp = i
		}
		if cleansDnfCache(seg) {
			lastClean = i
		}
	}
	return lastOp >= 0 && lastClean < lastOp
}

// isDnfOperation reports whether the segment invokes dnf/microdnf with modifying subcommands.
func isDnfOperation(seg []string) bool {
	if len(seg) < 2 {
		return false
	}
	if seg[0] != "dnf" && seg[0] != "microdnf" {
		return false
	}
	for _, t := range seg[1:] {
		switch t {
		case "install", "upgrade", "update", "groupinstall", "groupupdate", "distrosync", "autoremove", "remove":
			return true
		}
	}
	return false
}

// cleansDnfCache reports whether the segment removes dnf metadata/cache.
func cleansDnfCache(seg []string) bool {
	if len(seg) == 0 {
		return false
	}
	switch seg[0] {
	case "dnf", "microdnf":
		for i := 1; i < len(seg); i++ {
			if seg[i] == "clean" && i+1 < len(seg) && seg[i+1] == "all" {
				return true
			}
		}
		return false
	case "rm":
		joined := strings.Join(seg, " ")
		if !strings.Contains(joined, "/var/cache/dnf") {
			return false
		}
		flags := strings.Join(seg[1:], " ")
		return strings.Contains(flags, "-rf") || strings.Contains(flags, "-fr") || strings.Contains(flags, "-r")
	case "find":
		joined := strings.Join(seg, " ")
		if !strings.Contains(joined, "/var/cache/dnf") {
			return false
		}
		return strings.Contains(joined, "-delete")
	default:
		return false
	}
}
