// file: internal/rules/DL3001.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// noIrrelevantCommands flags usage of meaningless commands inside containers.
type noIrrelevantCommands struct{}

// NewNoIrrelevantCommands constructs the rule.
func NewNoIrrelevantCommands() engine.Rule { return noIrrelevantCommands{} }

// ID returns the rule identifier.
func (noIrrelevantCommands) ID() string { return "DL3001" }

// Check inspects RUN instructions for disallowed commands.
func (noIrrelevantCommands) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	invalid := map[string]struct{}{
		"ssh": {}, "vim": {}, "shutdown": {}, "service": {}, "ps": {},
		"free": {}, "top": {}, "kill": {}, "mount": {}, "ifconfig": {},
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		cmds := extractCommands(n)
		for _, c := range cmds {
			if _, bad := invalid[c]; bad {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3001",
					Message: "For some bash commands it makes no sense running them in a Docker container like `ssh`, `vim`, `shutdown`, `service`, `ps`, `free`, `top`, `kill`, `mount`, `ifconfig`.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// extractCommands returns command names invoked in a RUN instruction.
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
