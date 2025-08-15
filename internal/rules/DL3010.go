package rules

/*
 * file: internal/rules/DL3010.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// useADDForArchives encourages use of ADD for local archive extraction.
type useADDForArchives struct{}

// NewUseADDForArchives constructs the rule.
func NewUseADDForArchives() engine.Rule { return useADDForArchives{} }

// ID returns the rule identifier.
func (useADDForArchives) ID() string { return "DL3010" }

// Check examines COPY instructions that copy local tar archives to directories.
func (useADDForArchives) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "copy") {
			continue
		}
		if hasFromFlag(n.Flags) {
			continue
		}
		tokens := collectArgs(n)
		if len(tokens) < 2 {
			continue
		}
		dest := tokens[len(tokens)-1]
		if !strings.HasSuffix(dest, "/") {
			continue
		}
		for _, src := range tokens[:len(tokens)-1] {
			if isTarArchive(src) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3010",
					Message: "Instead of 'COPY " + src + " " + dest + "', use 'ADD " + src + " " + dest + "' to auto-extract.",
					Line:    n.StartLine,
				})
				break
			}
		}
	}
	return findings, nil
}

// hasFromFlag reports whether any flag is --from.
func hasFromFlag(flags []string) bool {
	for _, f := range flags {
		if strings.HasPrefix(strings.ToLower(f), "--from=") {
			return true
		}
	}
	return false
}

// collectArgs returns the arguments of an instruction as a slice of tokens.
func collectArgs(n *parser.Node) []string {
	var args []string
	for arg := n.Next; arg != nil; arg = arg.Next {
		args = append(args, strings.Trim(arg.Value, "\"'"))
	}
	return args
}

// isTarArchive reports whether the provided path appears to be a tar archive.
func isTarArchive(p string) bool {
	lp := strings.ToLower(p)
	switch {
	case strings.HasSuffix(lp, ".tar"),
		strings.HasSuffix(lp, ".tar.gz"),
		strings.HasSuffix(lp, ".tgz"),
		strings.HasSuffix(lp, ".tar.bz2"),
		strings.HasSuffix(lp, ".tbz"),
		strings.HasSuffix(lp, ".tar.xz"),
		strings.HasSuffix(lp, ".txz"):
		return true
	}
	return false
}
