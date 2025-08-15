package rules

/*
 * file: internal/rules/DL3020.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// useCopyInsteadOfAdd advises using COPY for local files and folders.
type useCopyInsteadOfAdd struct{}

// NewUseCopyInsteadOfAdd constructs the rule.
func NewUseCopyInsteadOfAdd() engine.Rule { return useCopyInsteadOfAdd{} }

// ID returns the rule identifier.
func (useCopyInsteadOfAdd) ID() string { return "DL3020" }

// Check inspects ADD instructions for local sources that are neither URLs nor archives.
func (useCopyInsteadOfAdd) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "add") {
			continue
		}
		tokens := collectArgs(n)
		if len(tokens) < 2 {
			continue
		}
		sources := tokens[:len(tokens)-1]
		for _, src := range sources {
			if !isURL(src) && !isArchive(src) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3020",
					Message: "Use COPY instead of ADD for files and folders",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// isURL reports whether the path appears to be a remote URL.
func isURL(p string) bool {
	lp := strings.ToLower(p)
	return strings.HasPrefix(lp, "http://") || strings.HasPrefix(lp, "https://")
}

var archiveFileExtensions = []string{
	".tar",
	".z",
	".bz2",
	".gz",
	".lz",
	".lzma",
	".tz",
	".tb2",
	".tbz",
	".tbz2",
	".tgz",
	".tlz",
	".tpz",
	".txz",
	".xz",
}

// isArchive reports whether the path looks like a compressed archive.
func isArchive(p string) bool {
	lp := strings.ToLower(p)
	for _, ext := range archiveFileExtensions {
		if strings.HasSuffix(lp, ext) {
			return true
		}
	}
	return false
}
