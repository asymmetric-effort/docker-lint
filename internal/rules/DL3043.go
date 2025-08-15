package rules

/*
 * file: internal/rules/DL3043.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"
	"unicode"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// requireOSVersionTag ensures OS base images specify explicit version tags.
type requireOSVersionTag struct{}

// NewRequireOSVersionTag constructs the rule.
func NewRequireOSVersionTag() engine.Rule { return requireOSVersionTag{} }

// ID returns the rule identifier.
func (requireOSVersionTag) ID() string { return "DL3043" }

// Check validates that OS images include explicit numeric version tags.
func (requireOSVersionTag) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil {
		return findings, nil
	}
	aliases := map[string]struct{}{}
	for _, s := range d.Stages {
		from := s.From
		if from == "" {
			if s.Name != "" {
				aliases[strings.ToLower(s.Name)] = struct{}{}
			}
			continue
		}
		if _, ok := aliases[strings.ToLower(from)]; ok {
			if s.Name != "" {
				aliases[strings.ToLower(s.Name)] = struct{}{}
			}
			continue
		}
		if needsOSVersionTag(from) {
			line := 0
			if s.Node != nil {
				line = s.Node.StartLine
			}
			findings = append(findings, engine.Finding{
				RuleID:  "DL3043",
				Message: "Specify OS version tag in FROM image",
				Line:    line,
			})
		}
		if s.Name != "" {
			aliases[strings.ToLower(s.Name)] = struct{}{}
		}
	}
	return findings, nil
}

// needsOSVersionTag reports if the image is OS-based and lacks a pinned version.
func needsOSVersionTag(image string) bool {
	lower := strings.ToLower(image)
	if lower == "scratch" {
		return false
	}
	if strings.HasPrefix(lower, "$") {
		return false
	}
	if !isOSImage(lower) {
		return false
	}
	name := lower
	if idx := strings.Index(name, "@"); idx != -1 {
		name = name[:idx]
	}
	parts := strings.Split(name, ":")
	if len(parts) == 1 {
		return true
	}
	tag := parts[len(parts)-1]
	floating := map[string]struct{}{"latest": {}, "stable": {}, "edge": {}, "rolling": {}}
	if _, ok := floating[tag]; ok {
		return true
	}
	for _, r := range tag {
		if unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// isOSImage determines if the image references a known OS base.
func isOSImage(image string) bool {
	name := image
	if idx := strings.Index(name, "@"); idx != -1 {
		name = name[:idx]
	}
	if idx := strings.Index(name, ":"); idx != -1 {
		name = name[:idx]
	}
	parts := strings.Split(name, "/")
	osNames := map[string]struct{}{
		"alpine": {}, "ubuntu": {}, "debian": {}, "fedora": {}, "centos": {},
		"rockylinux": {}, "rocky": {}, "almalinux": {}, "opensuse": {}, "suse": {},
		"archlinux": {}, "arch": {}, "amazonlinux": {}, "oraclelinux": {},
	}
	for _, p := range parts {
		if _, ok := osNames[p]; ok {
			return true
		}
	}
	return false
}
