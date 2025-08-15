package rules

/*
 * file: internal/rules/DL3047.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// apkCacheCleanup ensures APK cache is cleared or avoided after installation.
type apkCacheCleanup struct{}

// NewApkCacheCleanup constructs the rule.
func NewApkCacheCleanup() engine.Rule { return apkCacheCleanup{} }

// ID returns the rule identifier.
func (apkCacheCleanup) ID() string { return "DL3047" }

// Check verifies apk add usage includes --no-cache or cache removal.
func (apkCacheCleanup) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := lowerSegments(splitRunSegments(n))
		if needsApkCacheCleanup(segments) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3047",
				Message: "Clean apk cache after installing packages. Use 'apk add --no-cache' or remove /var/cache/apk/* in the same RUN.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// needsApkCacheCleanup reports whether an apk add occurred without --no-cache and without subsequent cache removal.
func needsApkCacheCleanup(segments [][]string) bool {
	lastInstall := -1
	installHasNoCache := false
	lastCleanup := -1
	for i, seg := range segments {
		if isApkAdd(seg) {
			lastInstall = i
			installHasNoCache = hasNoCache(seg)
		}
		if removesApkCache(seg) {
			lastCleanup = i
		}
	}
	return lastInstall >= 0 && !installHasNoCache && lastCleanup < lastInstall
}

// removesApkCache reports whether segment removes /var/cache/apk files.
func removesApkCache(seg []string) bool {
	if len(seg) == 0 {
		return false
	}
	joined := strings.Join(seg, " ")
	if !strings.Contains(joined, "/var/cache/apk") {
		return false
	}
	switch seg[0] {
	case "rm":
		flags := strings.Join(seg[1:], " ")
		return strings.Contains(flags, "-rf") || strings.Contains(flags, "-fr") || strings.Contains(flags, "-r")
	case "find":
		return strings.Contains(joined, "-delete")
	default:
		return false
	}
}
