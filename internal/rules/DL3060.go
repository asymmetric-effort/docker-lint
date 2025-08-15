package rules

/*
 * file: internal/rules/DL3060.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/google/shlex"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// yarnCacheClean ensures yarn cache is cleaned after yarn install.
type yarnCacheClean struct{}

// NewYarnCacheClean constructs the rule.
func NewYarnCacheClean() engine.Rule { return yarnCacheClean{} }

// ID returns the rule identifier.
func (yarnCacheClean) ID() string { return "DL3060" }

// Check warns when `yarn install` is used without subsequent `yarn cache clean`.
func (yarnCacheClean) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		if hasCacheMount(n.Flags) {
			continue
		}
		if n.Next == nil {
			continue
		}
		tokens, err := shlex.Split(n.Next.Value)
		if err != nil {
			continue
		}
		segments := splitByConnectors(tokens)
		install := false
		clean := false
		for _, seg := range segments {
			if isYarnInstall(seg) {
				install = true
			}
			if isYarnCacheClean(seg) {
				clean = true
			}
		}
		if install && !clean {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3060",
				Message: "`yarn cache clean` missing after `yarn install` was run.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// hasCacheMount reports whether a cache mount is present.
func hasCacheMount(flags []string) bool {
	for _, f := range flags {
		lf := strings.ToLower(f)
		if !strings.HasPrefix(lf, "--mount=") {
			continue
		}
		opts := strings.Split(strings.TrimPrefix(lf, "--mount="), ",")
		for _, o := range opts {
			if strings.HasPrefix(o, "type=") && strings.TrimPrefix(o, "type=") == "cache" {
				return true
			}
		}
	}
	return false
}

// isYarnInstall reports whether tokens represent `yarn install`.
func isYarnInstall(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}
	if strings.ToLower(tokens[0]) != "yarn" {
		return false
	}
	return strings.ToLower(tokens[1]) == "install"
}

// isYarnCacheClean reports whether tokens represent `yarn cache clean`.
func isYarnCacheClean(tokens []string) bool {
	if len(tokens) < 3 {
		return false
	}
	if strings.ToLower(tokens[0]) != "yarn" {
		return false
	}
	if strings.ToLower(tokens[1]) != "cache" {
		return false
	}
	return strings.ToLower(tokens[2]) == "clean"
}
