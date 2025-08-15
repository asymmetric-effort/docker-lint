package rules

/*
 * file: internal/rules/DL3055.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

var digestPattern = regexp.MustCompile(`@sha256:[0-9a-f]{64}$`)

// stageDigestPinned ensures configured stages pin images by digest.
type stageDigestPinned struct{ required map[string]struct{} }

// NewStageDigestPinned constructs the rule.
func NewStageDigestPinned(stages []string) engine.Rule {
	req := make(map[string]struct{}, len(stages))
	for _, s := range stages {
		req[strings.ToLower(s)] = struct{}{}
	}
	return &stageDigestPinned{required: req}
}

// ID returns the rule identifier.
func (stageDigestPinned) ID() string { return "DL3055" }

// Check verifies required stages use digest-pinned images.
func (r *stageDigestPinned) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	if len(r.required) == 0 {
		return findings, nil
	}
	for _, st := range d.Stages {
		if _, ok := r.required[strings.ToLower(st.Name)]; !ok {
			continue
		}
		img := strings.ToLower(st.From)
		if !digestPattern.MatchString(img) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3055",
				Message: "Stage \"" + st.Name + "\" image is not pinned by digest.",
				Line:    st.Node.StartLine,
			})
		}
	}
	return findings, nil
}
