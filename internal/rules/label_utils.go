package rules

/*
 * file: internal/rules/label_utils.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// LabelType represents the expected format for a label value.
type LabelType int

const (
	LabelTypeString LabelType = iota
	LabelTypeURL
	LabelTypeRFC3339
	LabelTypeSPDX
	LabelTypeGitHash
	LabelTypeSemVer
	LabelTypeEmail
)

// LabelSchema defines required labels and their expected types.
type LabelSchema map[string]LabelType

// labelPair holds a key-value label entry.
type labelPair struct{ Key, Value string }

// collectLabelPairs extracts key-value pairs from a LABEL instruction.
func collectLabelPairs(n *parser.Node) []labelPair {
	var pairs []labelPair
	if n == nil {
		return pairs
	}
	var tokens []string
	for tok := n.Next; tok != nil; tok = tok.Next {
		tokens = append(tokens, strings.Trim(tok.Value, "\"'"))
	}
	for i := 0; i+2 < len(tokens); i += 3 {
		pairs = append(pairs, labelPair{Key: tokens[i], Value: tokens[i+1]})
	}
	return pairs
}

// inSchema reports whether a key exists in the schema.
func inSchema(schema LabelSchema, key string) bool {
	_, ok := schema[key]
	return ok
}
