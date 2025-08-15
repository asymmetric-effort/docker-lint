package rules

/*
 * file: internal/rules/DL3009.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// aptListsCleanup ensures apt installs remove package lists in the same layer.
type aptListsCleanup struct{}

// NewAptListsCleanup constructs the rule.
func NewAptListsCleanup() engine.Rule { return aptListsCleanup{} }

// ID returns the rule identifier.
func (aptListsCleanup) ID() string { return "DL3009" }

// Check scans RUN instructions for apt installs lacking cleanup.
func (aptListsCleanup) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := lowerSegments(splitRunSegments(n))
		if needsAptListsCleanup(segments) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3009",
				Message: "Delete the APT lists after installing packages.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// needsAptListsCleanup reports whether an apt install occurred without subsequent cleanup.
func needsAptListsCleanup(segments [][]string) bool {
	lastInstall := -1
	lastCleanup := -1
	for i, seg := range segments {
		if isAptInstall(seg) {
			lastInstall = i
		}
		if removesAptLists(seg) {
			lastCleanup = i
		}
	}
	return lastInstall >= 0 && lastCleanup < lastInstall
}

// isAptInstall reports whether the segment invokes apt-get install or apt install.
func isAptInstall(seg []string) bool {
	if len(seg) < 2 {
		return false
	}
	if seg[0] != "apt-get" && seg[0] != "apt" {
		return false
	}
	for _, t := range seg[1:] {
		if t == "install" {
			return true
		}
	}
	return false
}

// removesAptLists reports whether the segment deletes the apt lists directory.
func removesAptLists(seg []string) bool {
	if len(seg) == 0 {
		return false
	}
	joined := strings.Join(seg, " ")
	if !strings.Contains(joined, "/var/lib/apt/lists") {
		return false
	}
	switch seg[0] {
	case "rm":
		flags := strings.Join(seg[1:], " ")
		return strings.Contains(flags, "-rf") || strings.Contains(flags, "-fr") || strings.Contains(flags, "-r")
	case "find":
		return strings.Contains(joined, "-delete")
	default:
		return false
	}
}
