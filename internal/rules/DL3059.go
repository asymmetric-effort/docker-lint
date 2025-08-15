package rules

/*
 * file: internal/rules/DL3059.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"sort"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// consecutiveRun detects multiple consecutive RUN instructions with few commands.
type consecutiveRun struct {
	prevFlags string
	prevCount int
	seen      bool
}

// NewConsecutiveRun constructs the rule.
func NewConsecutiveRun() engine.Rule { return &consecutiveRun{} }

// ID returns the rule identifier.
func (consecutiveRun) ID() string { return "DL3059" }

// Check flags consecutive simple RUN instructions for consolidation.
func (r *consecutiveRun) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	reset := func() { r.seen = false }
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			if strings.EqualFold(n.Value, "#") {
				continue
			}
			reset()
			continue
		}
		count := countRunCommands(n)
		flags := canonicalFlags(n.Flags)
		if r.seen && r.prevFlags == flags && r.prevCount <= 2 && count <= 2 {
			findings = append(findings, engine.Finding{RuleID: "DL3059", Message: "Multiple consecutive `RUN` instructions. Consider consolidation.", Line: n.StartLine})
		}
		r.prevFlags = flags
		r.prevCount = count
		r.seen = true
	}
	return findings, nil
}

// canonicalFlags normalizes run flags for comparison.
func canonicalFlags(flags []string) string {
	cp := append([]string(nil), flags...)
	sort.Strings(cp)
	return strings.Join(cp, " ")
}

// countRunCommands returns the number of commands in a RUN instruction.
func countRunCommands(n *parser.Node) int {
	segs := splitRunSegments(n)
	return len(segs)
}
