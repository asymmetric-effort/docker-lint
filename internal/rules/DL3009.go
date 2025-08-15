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
	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
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
		segments := splitRunSegments(n)
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

// splitRunSegments splits a RUN instruction into command segments separated by shell connectors.
func splitRunSegments(n *parser.Node) [][]string {
	if n == nil || n.Next == nil {
		return nil
	}
	if n.Attributes != nil && n.Attributes["json"] {
		return [][]string{{strings.ToLower(n.Next.Value)}}
	}
	tokens, err := shlex.Split(n.Next.Value)
	if err != nil {
		return nil
	}
	var segs [][]string
	var cur []string
	for _, tok := range tokens {
		switch tok {
		case "&&", "||", ";", "|":
			if len(cur) > 0 {
				segs = append(segs, lowerSlice(cur))
				cur = nil
			}
		default:
			cur = append(cur, tok)
		}
	}
	if len(cur) > 0 {
		segs = append(segs, lowerSlice(cur))
	}
	return segs
}

// lowerSlice returns a new slice with all elements lowercased.
func lowerSlice(in []string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = strings.ToLower(s)
	}
	return out
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
