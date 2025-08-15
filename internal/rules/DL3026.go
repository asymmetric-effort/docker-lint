// file: internal/rules/DL3026.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// allowedRegistry enforces FROM images use only allowed registries.
type allowedRegistry struct {
	allowed []string
}

// NewAllowedRegistry constructs the rule.
func NewAllowedRegistry(allowed []string) engine.Rule { return &allowedRegistry{allowed: allowed} }

// ID returns the rule identifier.
func (allowedRegistry) ID() string { return "DL3026" }

// Check validates registry usage in FROM instructions.
func (r *allowedRegistry) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	aliases := map[string]struct{}{}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "from") {
			continue
		}
		image := ""
		if n.Next != nil {
			image = n.Next.Value
		}
		alias := strings.ToLower(stageAlias(n))
		if alias != "" {
			aliases[alias] = struct{}{}
		}
		if _, ok := aliases[strings.ToLower(image)]; ok {
			continue // referencing previous stage
		}
		if len(r.allowed) == 0 {
			continue
		}
		if !r.isAllowed(image) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3026",
				Message: "Use only an allowed registry in the FROM image",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

// isAllowed reports if the image registry is permitted.
func (r *allowedRegistry) isAllowed(image string) bool {
	registry := parseRegistry(image)
	if registry == "" {
		if strings.EqualFold(image, "scratch") {
			return true
		}
		return r.match("docker.io") || r.match("hub.docker.com")
	}
	return r.match(registry)
}

// match checks allowed patterns against registry.
func (r *allowedRegistry) match(registry string) bool {
	for _, a := range r.allowed {
		if matchRegistryPattern(a, registry) {
			return true
		}
	}
	return false
}

// parseRegistry extracts the registry part of an image reference.
func parseRegistry(image string) string {
	parts := strings.Split(image, "/")
	if len(parts) > 1 && (strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") || parts[0] == "localhost") {
		return parts[0]
	}
	return ""
}

// matchRegistryPattern matches registry against pattern with '*' prefix or suffix.
func matchRegistryPattern(pattern, registry string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(registry, pattern[1:])
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(registry, pattern[:len(pattern)-1])
	}
	return registry == pattern
}
