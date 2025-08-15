package rules

/*
 * file: internal/rules/DL3014.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/google/shlex"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// aptGetYes enforces non-interactive apt-get install invocations.
type aptGetYes struct{}

// NewAptGetYes constructs the rule.
func NewAptGetYes() engine.Rule { return aptGetYes{} }

// ID returns the rule identifier.
func (aptGetYes) ID() string { return "DL3014" }

// Check scans RUN instructions for apt-get install lacking -y or equivalent.
func (aptGetYes) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		tokens, err := shlex.Split(n.Next.Value)
		if err != nil {
			continue
		}
		segments := splitByConnectors(tokens)
		for _, seg := range segments {
			if isAptGetInstall(seg) && !hasYesOption(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3014",
					Message: "Use the -y switch to avoid manual input apt-get -y install <package>",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// splitByConnectors divides tokens into command segments.
func splitByConnectors(tokens []string) [][]string {
	var segments [][]string
	var current []string
	for _, t := range tokens {
		switch t {
		case "&&", "||", "|", ";":
			if len(current) > 0 {
				segments = append(segments, current)
				current = nil
			}
		default:
			current = append(current, strings.ToLower(t))
		}
	}
	if len(current) > 0 {
		segments = append(segments, current)
	}
	return segments
}

// isAptGetInstall reports whether tokens represent `apt-get install`.
func isAptGetInstall(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokens[0] != "apt-get" {
		return false
	}
	for _, t := range tokens[1:] {
		if t == "install" {
			return true
		}
	}
	return false
}

// hasYesOption detects non-interactive flags for apt-get.
func hasYesOption(tokens []string) bool {
	yesFlags := map[string]struct{}{
		"-y":           {},
		"--yes":        {},
		"--assume-yes": {},
		"-qq":          {},
	}
	qCount := 0
	quietCount := 0
	for _, t := range tokens[1:] {
		if _, ok := yesFlags[t]; ok {
			return true
		}
		if t == "-q=2" || t == "--quiet=2" {
			return true
		}
		if t == "-q" {
			qCount++
		}
		if t == "--quiet" {
			quietCount++
		}
	}
	return qCount >= 2 || quietCount >= 2
}
