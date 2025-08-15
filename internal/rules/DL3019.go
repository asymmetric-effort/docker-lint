package rules

/*
 * file: internal/rules/DL3019.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/google/shlex"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// apkNoCache enforces use of --no-cache for apk add unless cache mount is present.
type apkNoCache struct{}

// NewApkNoCache constructs the rule.
func NewApkNoCache() engine.Rule { return apkNoCache{} }

// ID returns the rule identifier.
func (apkNoCache) ID() string { return "DL3019" }

// Check examines RUN instructions for apk add missing --no-cache and without cache mount.
func (apkNoCache) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		if hasApkCacheMount(n.Flags) {
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
		for _, seg := range segments {
			if isApkAdd(seg) && !hasNoCache(seg) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3019",
					Message: "Use the `--no-cache` switch to avoid the need to use `--update` and remove `/var/cache/apk/*` when done installing packages",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// hasApkCacheMount reports whether a cache mount targets /var/cache/apk.
func hasApkCacheMount(flags []string) bool {
	for _, f := range flags {
		lf := strings.ToLower(f)
		if !strings.HasPrefix(lf, "--mount=") {
			continue
		}
		opts := strings.Split(strings.TrimPrefix(lf, "--mount="), ",")
		typeCache := false
		targetMatch := false
		for _, o := range opts {
			if strings.HasPrefix(o, "type=") && strings.TrimPrefix(o, "type=") == "cache" {
				typeCache = true
			}
			if strings.HasPrefix(o, "target=") || strings.HasPrefix(o, "dst=") || strings.HasPrefix(o, "destination=") {
				p := strings.TrimPrefix(o, "target=")
				p = strings.TrimPrefix(p, "dst=")
				p = strings.TrimPrefix(p, "destination=")
				p = strings.TrimRight(p, "/")
				if p == "/var/cache/apk" {
					targetMatch = true
				}
			}
		}
		if typeCache && targetMatch {
			return true
		}
	}
	return false
}

// isApkAdd reports whether tokens start with `apk` and include `add`.
func isApkAdd(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}
	if strings.ToLower(tokens[0]) != "apk" {
		return false
	}
	for _, t := range tokens[1:] {
		if strings.ToLower(t) == "add" {
			return true
		}
	}
	return false
}

// hasNoCache checks for --no-cache flag.
func hasNoCache(tokens []string) bool {
	for _, t := range tokens[1:] {
		if strings.ToLower(t) == "--no-cache" {
			return true
		}
	}
	return false
}
