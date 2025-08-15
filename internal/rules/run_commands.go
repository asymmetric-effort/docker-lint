// file: internal/rules/run_commands.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"

	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// extractCommands returns command names invoked in a RUN instruction.
//
// extractCommands inspects the RUN node and returns the list of command
// names, respecting shell parsing and handling JSON-form RUN variants.
func extractCommands(n *parser.Node) []string {
	if n == nil || n.Next == nil {
		return nil
	}
	if n.Attributes != nil && n.Attributes["json"] {
		return []string{strings.ToLower(n.Next.Value)}
	}
	tokens, err := shlex.Split(n.Next.Value)
	if err != nil {
		return nil
	}
	return commandNames(tokens)
}

// commandNames identifies command boundaries within shell tokens.
//
// commandNames returns the ordered list of commands invoked within the
// tokenized shell statement, tracking boundaries across common shell
// connectors like && and ||.
func commandNames(tokens []string) []string {
	var cmds []string
	expect := true
	for _, tok := range tokens {
		if expect {
			cmds = append(cmds, strings.ToLower(tok))
			expect = false
			continue
		}
		switch tok {
		case "&&", "||", "|", ";":
			expect = true
		}
	}
	return cmds
}

// splitRunSegments tokenizes a RUN instruction and splits it into command segments.
// It handles both shell-form and JSON-form RUN instructions.
func splitRunSegments(n *parser.Node) [][]string {
	if n == nil || n.Next == nil {
		return nil
	}
	var tokens []string
	if n.Attributes != nil && n.Attributes["json"] {
		for tok := n.Next; tok != nil; tok = tok.Next {
			tokens = append(tokens, tok.Value)
		}
	} else {
		t, err := shlex.Split(n.Next.Value)
		if err != nil {
			return nil
		}
		tokens = t
	}
	var segments [][]string
	var current []string
	for _, tok := range tokens {
		switch tok {
		case "&&", "||", "|", ";":
			if len(current) > 0 {
				segments = append(segments, current)
				current = nil
			}
		default:
			current = append(current, tok)
		}
	}
	if len(current) > 0 {
		segments = append(segments, current)
	}
	return segments
}

// lowerSegments returns a lowercase copy of each segment.
func lowerSegments(segs [][]string) [][]string {
	out := make([][]string, len(segs))
	for i, seg := range segs {
		outSeg := make([]string, len(seg))
		for j, s := range seg {
			outSeg[j] = strings.ToLower(s)
		}
		out[i] = outSeg
	}
	return out
}
