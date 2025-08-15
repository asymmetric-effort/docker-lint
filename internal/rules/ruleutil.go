// file: internal/rules/ruleutil.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

package rules

import (
	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"strings"
)

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
