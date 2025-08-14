// file: internal/engine/engine.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package engine

import (
	"context"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// Finding represents a lint result.
//
// Finding identifies a rule violation with a message and line number.
type Finding struct {
	RuleID  string `json:"rule"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
}

// Rule defines the interface for lint rules.
//
// Check evaluates the document and returns any findings.
type Rule interface {
	ID() string
	Check(ctx context.Context, d *ir.Document) ([]Finding, error)
}

// Registry stores and executes lint rules.
//
// Registry allows registration of rules and running them over a document.
type Registry struct {
	rules []Rule
}

// NewRegistry creates an empty rule registry.
func NewRegistry() *Registry { return &Registry{} }

// Register adds a rule to the registry.
func (r *Registry) Register(rule Rule) { r.rules = append(r.rules, rule) }

// Run executes all registered rules against the document.
func (r *Registry) Run(ctx context.Context, d *ir.Document) ([]Finding, error) {
	var all []Finding
	for _, rl := range r.rules {
		f, err := rl.Check(ctx, d)
		if err != nil {
			return nil, err
		}
		all = append(all, f...)
	}
	return all, nil
}
