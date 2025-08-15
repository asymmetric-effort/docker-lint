package rules

/*
 * file: internal/rules/DL3013.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/google/shlex"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// pinPipVersions enforces version pinning for pip installations.
type pinPipVersions struct{}

// NewPinPipVersions constructs the rule.
func NewPinPipVersions() engine.Rule { return pinPipVersions{} }

// ID returns the rule identifier.
func (pinPipVersions) ID() string { return "DL3013" }

// Check inspects RUN instructions for unpinned pip installs.
func (pinPipVersions) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") || n.Next == nil {
			continue
		}
		tokens, err := shlex.Split(n.Next.Value)
		if err != nil {
			continue
		}
		cmds := splitRunCommands(tokens)
		for _, cmd := range cmds {
			if violatesPipPin(cmd) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3013",
					Message: "Pin versions in pip. Instead of `pip install <package>` use `pip install <package>==<version>` or `pip install --requirement <requirements file>`",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// splitRunCommands divides tokens into individual commands based on shell connectors.
func splitRunCommands(tokens []string) [][]string {
	connectors := map[string]struct{}{"&&": {}, "||": {}, "|": {}, ";": {}}
	var result [][]string
	var current []string
	for _, tok := range tokens {
		if _, ok := connectors[tok]; ok {
			if len(current) > 0 {
				result = append(result, current)
				current = nil
			}
			continue
		}
		current = append(current, tok)
	}
	if len(current) > 0 {
		result = append(result, current)
	}
	return result
}

// violatesPipPin reports whether a pip install command lacks version pinning.
func violatesPipPin(cmd []string) bool {
	start, ok := pipInstallIndex(cmd)
	if !ok {
		return false
	}
	hasConstraint := false
	requirement := false
	var pkgs []string
	flagsWithArg := map[string]struct{}{
		"abi": {}, "b": {}, "build": {}, "e": {}, "editable": {}, "extra-index-url": {},
		"f": {}, "find-links": {}, "i": {}, "index-url": {}, "implementation": {},
		"no-binary": {}, "only-binary": {}, "platform": {}, "prefix": {}, "progress-bar": {},
		"proxy": {}, "python-version": {}, "root": {}, "src": {}, "t": {}, "target": {},
		"trusted-host": {}, "upgrade-strategy": {},
	}
	for i := start; i < len(cmd); i++ {
		tok := cmd[i]
		if strings.HasPrefix(tok, "-") {
			flag := strings.TrimLeft(tok, "-")
			switch flag {
			case "r", "requirement":
				requirement = true
				i++
				break
			case "c", "constraint":
				hasConstraint = true
				i++
				break
			default:
				if _, ok := flagsWithArg[flag]; ok {
					i++
				}
			}
			continue
		}
		if tok == "." {
			requirement = true
			break
		}
		pkgs = append(pkgs, tok)
	}
	if requirement || hasConstraint || len(pkgs) == 0 {
		return false
	}
	for _, p := range pkgs {
		if !versionFixed(p) {
			return true
		}
	}
	return false
}

// pipInstallIndex identifies pip install command and returns index after install token.
func pipInstallIndex(cmd []string) (int, bool) {
	if len(cmd) >= 2 && isPip(cmd[0]) && cmd[1] == "install" {
		return 2, true
	}
	if len(cmd) >= 4 && strings.HasPrefix(cmd[0], "python") && cmd[1] == "-m" && isPip(cmd[2]) && cmd[3] == "install" {
		return 4, true
	}
	return 0, false
}

// isPip reports whether the token refers to pip.
func isPip(tok string) bool {
	return strings.HasPrefix(tok, "pip")
}

// versionFixed reports whether a package token pins its version.
func versionFixed(pkg string) bool {
	if strings.Contains(pkg, "@") {
		return true
	}
	symbols := []string{"==", ">=", "<=", ">", "<", "!=", "~=", "==="}
	for _, s := range symbols {
		if strings.Contains(pkg, s) {
			return true
		}
	}
	if strings.HasSuffix(pkg, ".whl") || strings.HasSuffix(pkg, ".tar.gz") {
		return true
	}
	if strings.Contains(pkg, "/") {
		return true
	}
	return false
}
