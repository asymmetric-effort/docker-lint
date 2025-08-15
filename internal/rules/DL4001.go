// file: internal/rules/DL4001.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// exclusiveCurlWget ensures only one of curl or wget is used per stage.
type exclusiveCurlWget struct{}

// NewExclusiveCurlWget constructs the rule.
func NewExclusiveCurlWget() engine.Rule { return exclusiveCurlWget{} }

// ID returns the rule identifier.
func (exclusiveCurlWget) ID() string { return "DL4001" }

// Check scans RUN instructions for mixed use of curl and wget.
func (exclusiveCurlWget) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	seenCurl := false
	seenWget := false
	for _, n := range d.AST.Children {
		if strings.EqualFold(n.Value, "from") {
			seenCurl = false
			seenWget = false
			continue
		}
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		cmds := extractCommands(n)
		usesCurl := false
		usesWget := false
		for _, cmd := range cmds {
			switch cmd {
			case "curl":
				usesCurl = true
			case "wget":
				usesWget = true
			}
		}
		if usesCurl && usesWget || (usesCurl && seenWget) || (usesWget && seenCurl) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL4001",
				Message: "Either use Wget or Curl but not both",
				Line:    n.StartLine,
			})
		}
		if usesCurl {
			seenCurl = true
		}
		if usesWget {
			seenWget = true
		}
	}
	return findings, nil
}
