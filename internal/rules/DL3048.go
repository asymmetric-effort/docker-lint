package rules

/*
 * file: internal/rules/DL3048.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

var reservedLabelNS = regexp.MustCompile(`^(?:com|io)\.docker\.|^org\.dockerproject\.`)

// labelKeyValid ensures LABEL keys follow Docker recommendations and avoid reserved namespaces.
type labelKeyValid struct{}

// NewLabelKeyValid constructs the rule.
func NewLabelKeyValid() engine.Rule { return labelKeyValid{} }

// ID returns the rule identifier.
func (labelKeyValid) ID() string { return "DL3048" }

// Check validates label keys for reserved namespaces and allowed characters.
func (labelKeyValid) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		switch {
		case strings.EqualFold(n.Value, "label"):
			findings = append(findings, validateLabelNode(n, n)...)
		case strings.EqualFold(n.Value, "onbuild") && n.Next != nil:
			for _, c := range n.Next.Children {
				if strings.EqualFold(c.Value, "label") {
					findings = append(findings, validateLabelNode(c, n)...)
				}
			}
		}
	}
	return findings, nil
}

// validateLabelNode reports invalid label keys from the given LABEL node.
func validateLabelNode(ln *parser.Node, lineNode *parser.Node) []engine.Finding {
	var findings []engine.Finding
	for _, p := range collectLabelPairs(ln) {
		if invalidLabelKey(p.Key) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3048",
				Message: "Label key `" + p.Key + "` is invalid. Use lower-case a–z, 0–9, '.' and '-' only; avoid reserved namespaces; no leading/trailing or repeated separators.",
				Line:    lineNode.StartLine,
			})
		}
	}
	return findings
}

// invalidLabelKey reports whether the key violates format restrictions or reserved namespaces.
func invalidLabelKey(key string) bool {
	if reservedLabelNS.MatchString(key) {
		return true
	}
	if key == "" {
		return true
	}
	if key[0] == '-' || key[0] == '.' || key[len(key)-1] == '-' || key[len(key)-1] == '.' {
		return true
	}
	if strings.Contains(key, "..") || strings.Contains(key, "--") {
		return true
	}
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '.' || r == '-':
		default:
			return true
		}
	}
	return false
}
