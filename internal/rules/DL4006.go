// file: internal/rules/DL4006.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// pipefailBeforePipe flags RUN instructions with pipes without preceding SHELL -o pipefail.
type pipefailBeforePipe struct{}

// NewPipefailBeforePipe constructs the rule.
func NewPipefailBeforePipe() engine.Rule { return pipefailBeforePipe{} }

// ID returns the rule identifier.
func (pipefailBeforePipe) ID() string { return "DL4006" }

// Check evaluates the Dockerfile for missing pipefail.
func (pipefailBeforePipe) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	pipefail := false
	nonPosix := []string{"pwsh", "powershell", "cmd"}
	valid := map[string]bool{"/bin/bash": true, "/bin/zsh": true, "/bin/ash": true, "bash": true, "zsh": true, "ash": true}
	for _, n := range d.AST.Children {
		switch strings.ToLower(n.Value) {
		case "from":
			pipefail = false
		case "shell":
			if isNonPosixShell(n, nonPosix) {
				pipefail = true
			} else {
				pipefail = hasPipefailOption(n, valid)
			}
		case "run":
			if !pipefail && runHasPipe(n) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL4006",
					Message: "Set the SHELL option -o pipefail before RUN with a pipe in it. If you are using /bin/sh in an alpine image or if your shell is symlinked to busybox then consider explicitly setting your SHELL to /bin/ash, or disable this check",
					Line:    n.StartLine,
				})
			}
		}
	}
	return findings, nil
}

// isNonPosixShell reports whether the shell is non-POSIX and thus exempt.
func isNonPosixShell(n *parser.Node, shells []string) bool {
	if n == nil || n.Next == nil {
		return false
	}
	sh := strings.ToLower(n.Next.Value)
	for _, s := range shells {
		if strings.HasPrefix(sh, s) {
			return true
		}
	}
	return false
}

// hasPipefailOption checks for -o pipefail in a SHELL instruction.
func hasPipefailOption(n *parser.Node, valid map[string]bool) bool {
	if n == nil || n.Next == nil {
		return false
	}
	sh := strings.ToLower(n.Next.Value)
	if !valid[sh] {
		return false
	}
	for t := n.Next.Next; t != nil; t = t.Next {
		if t.Value == "-o" && t.Next != nil && strings.EqualFold(t.Next.Value, "pipefail") {
			return true
		}
	}
	return false
}

// runHasPipe detects whether a RUN command contains a pipe.
func runHasPipe(n *parser.Node) bool {
	if n == nil || n.Next == nil {
		return false
	}
	if n.Attributes != nil && n.Attributes["json"] {
		for t := n.Next; t != nil; t = t.Next {
			if strings.Contains(t.Value, "|") {
				return true
			}
		}
		return false
	}
	return strings.Contains(n.Next.Value, "|")
}
