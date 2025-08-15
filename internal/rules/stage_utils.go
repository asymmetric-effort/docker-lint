// file: internal/rules/stage_utils.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// copyFromFlag extracts the value of --from from a COPY node.
func copyFromFlag(n *parser.Node) (string, bool) {
	for _, f := range n.Flags {
		if strings.HasPrefix(strings.ToLower(f), "--from=") {
			v := strings.TrimPrefix(f, "--from=")
			v = strings.Trim(v, "\"'")
			return v, true
		}
	}
	return "", false
}

// stageAlias returns the alias specified in a FROM instruction.
func stageAlias(n *parser.Node) string {
	if n == nil || n.Next == nil {
		return ""
	}
	for tok := n.Next.Next; tok != nil; tok = tok.Next {
		if strings.EqualFold(tok.Value, "as") && tok.Next != nil {
			return tok.Next.Value
		}
	}
	return ""
}
