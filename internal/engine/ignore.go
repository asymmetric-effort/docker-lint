package engine

/*
 * file: internal/engine/ignore.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

const ignoreDirective = "hadolint ignore="

// lineIgnores returns rule IDs to skip keyed by line number.
//
// lineIgnores scans the document's AST for `hadolint ignore=` pragmas and
// records which rules should be skipped for each instruction line.
func lineIgnores(d *ir.Document) map[int]map[string]struct{} {
	m := make(map[int]map[string]struct{})
	if d == nil || d.AST == nil {
		return m
	}
	for _, n := range d.AST.Children {
		line := n.StartLine
		for _, com := range n.PrevComment {
			addIgnores(m, line, parseIgnorePragma(com))
		}
		addIgnores(m, line, parseIgnorePragma(n.Original))
	}
	return m
}

// addIgnores merges rule IDs into the map for a given line.
func addIgnores(m map[int]map[string]struct{}, line int, ids []string) {
	if len(ids) == 0 {
		return
	}
	set, ok := m[line]
	if !ok {
		set = make(map[string]struct{})
		m[line] = set
	}
	for _, id := range ids {
		set[strings.ToUpper(id)] = struct{}{}
	}
}

// parseIgnorePragma extracts rule IDs from a comment or instruction, case-insensitively.
func parseIgnorePragma(s string) []string {
	lower := strings.ToLower(s)
	idx := strings.Index(lower, ignoreDirective)
	if idx == -1 {
		return nil
	}
	rest := lower[idx+len(ignoreDirective):]
	rest = strings.TrimSpace(rest)
	fields := strings.FieldsFunc(rest, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' })
	var ids []string
	for _, f := range fields {
		if trimmed := strings.TrimSpace(f); trimmed != "" {
			ids = append(ids, trimmed)
		}
	}
	return ids
}

// shouldSkip reports whether the finding should be skipped based on ignore pragmas.
func shouldSkip(ignores map[int]map[string]struct{}, f Finding) bool {
	ids, ok := ignores[f.Line]
	if !ok {
		return false
	}
	_, skip := ids[strings.ToUpper(f.RuleID)]
	return skip
}
