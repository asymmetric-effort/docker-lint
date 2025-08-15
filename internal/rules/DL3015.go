package rules

/*
 * file: internal/rules/DL3015.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// aptNoInstallRecommends ensures apt-get install uses --no-install-recommends.
type aptNoInstallRecommends struct{}

// NewAptNoInstallRecommends constructs the rule.
func NewAptNoInstallRecommends() engine.Rule { return aptNoInstallRecommends{} }

// ID returns the rule identifier.
func (aptNoInstallRecommends) ID() string { return "DL3015" }

// Check scans RUN instructions for apt-get install missing --no-install-recommends.
func (aptNoInstallRecommends) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		tokens := runTokens(n)
		cmds := splitTokens(tokens)
		for _, c := range cmds {
			if aptInstallMissingFlag(c) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3015",
					Message: "Avoid additional packages by specifying `--no-install-recommends`",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// runTokens returns shell tokens for a RUN instruction.
func runTokens(n *parser.Node) []string {
	if n == nil || n.Next == nil {
		return nil
	}
	if n.Attributes != nil && n.Attributes["json"] {
		var toks []string
		for c := n.Next; c != nil; c = c.Next {
			toks = append(toks, c.Value)
		}
		return toks
	}
	tokens, err := shlex.Split(n.Next.Value)
	if err != nil {
		return nil
	}
	return tokens
}

// splitTokens divides tokens into individual commands.
func splitTokens(tokens []string) [][]string {
	var cmds [][]string
	var cur []string
	for _, t := range tokens {
		switch t {
		case "&&", "||", ";", "|":
			if len(cur) > 0 {
				cmds = append(cmds, cur)
				cur = nil
			}
		default:
			cur = append(cur, t)
		}
	}
	if len(cur) > 0 {
		cmds = append(cmds, cur)
	}
	return cmds
}

// aptInstallMissingFlag reports apt-get install commands lacking no-install-recommends.
func aptInstallMissingFlag(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	if strings.ToLower(tokens[0]) != "apt-get" {
		return false
	}
	hasInstall := false
	hasFlag := false
	for _, t := range tokens[1:] {
		lt := strings.ToLower(t)
		if lt == "install" {
			hasInstall = true
			continue
		}
		if strings.Contains(lt, "no-install-recommends") || strings.Contains(lt, "apt::install-recommends=false") {
			hasFlag = true
		}
	}
	return hasInstall && !hasFlag
}
