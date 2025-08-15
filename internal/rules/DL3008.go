package rules

/*
 * file: internal/rules/DL3008.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// aptPin ensures packages installed with apt-get or apt are version pinned.
type aptPin struct{}

// NewAptPin constructs the rule.
func NewAptPin() engine.Rule { return aptPin{} }

// ID returns the rule identifier.
func (aptPin) ID() string { return "DL3008" }

// Check evaluates RUN instructions for unpinned apt installs.
func (aptPin) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") || n.Next == nil {
			continue
		}
		cmd := n.Next.Value
		if hasUnpinnedAptInstall(cmd) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3008",
				Message: "Pin versions in apt-get install. Instead of 'apt-get install <pkg>' use 'apt-get install <pkg>=<version>'.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

var splitter = regexp.MustCompile(`\s*(?:&&|\|\||;)\s*`)

func hasUnpinnedAptInstall(cmd string) bool {
	cmd = strings.ReplaceAll(cmd, "\\\n", " ")
	cmd = strings.ReplaceAll(cmd, "\n", " ")
	parts := splitter.Split(cmd, -1)
	for _, part := range parts {
		tokens := strings.Fields(part)
		for i := 0; i < len(tokens); i++ {
			tok := tokens[i]
			if tok == "apt-get" || tok == "apt" {
				for j := i + 1; j < len(tokens); j++ {
					t := tokens[j]
					if t == "install" {
						if unpinnedPackages(tokens[j+1:]) {
							return true
						}
						break
					}
					if !strings.HasPrefix(t, "-") {
						break
					}
				}
			}
		}
	}
	return false
}

func unpinnedPackages(args []string) bool {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		if !strings.Contains(a, "=") {
			return true
		}
		parts := strings.SplitN(a, "=", 2)
		if parts[1] == "" || strings.HasPrefix(parts[1], "-") {
			return true
		}
	}
	return false
}
