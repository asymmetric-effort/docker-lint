package rules

/*
 * file: internal/rules/DL3018.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"regexp"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// apkPin ensures packages installed with apk add are version pinned.
type apkPin struct{}

// NewApkPin constructs the rule.
func NewApkPin() engine.Rule { return apkPin{} }

// ID returns the rule identifier.
func (apkPin) ID() string { return "DL3018" }

// Check evaluates RUN instructions for unpinned apk adds.
func (apkPin) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") || n.Next == nil {
			continue
		}
		cmd := n.Next.Value
		if hasUnpinnedApkAdd(cmd) {
			findings = append(findings, engine.Finding{
				RuleID:  "DL3018",
				Message: "Pin versions in apk add. Instead of 'apk add <package>' use 'apk add <package>=<version>'.",
				Line:    n.StartLine,
			})
		}
	}
	return findings, nil
}

var apkSplitter = regexp.MustCompile(`\s*(?:&&|\|\||;)\s*`)

func hasUnpinnedApkAdd(cmd string) bool {
	cmd = strings.ReplaceAll(cmd, "\\\n", " ")
	cmd = strings.ReplaceAll(cmd, "\n", " ")
	parts := apkSplitter.Split(cmd, -1)
	for _, part := range parts {
		tokens := strings.Fields(part)
		for i := 0; i < len(tokens); i++ {
			if tokens[i] == "apk" {
				for j := i + 1; j < len(tokens); j++ {
					t := tokens[j]
					if t == "add" {
						if unpinnedApkPackages(tokens[j+1:]) {
							return true
						}
						break
					}
					if !strings.HasPrefix(t, "-") {
						break
					}
				}
			}
		}
	}
	return false
}

func unpinnedApkPackages(args []string) bool {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		if strings.HasSuffix(a, ".apk") {
			continue
		}
		if !strings.Contains(a, "=") {
			return true
		}
		parts := strings.SplitN(a, "=", 2)
		if parts[1] == "" || strings.HasPrefix(parts[1], "-") {
			return true
		}
	}
	return false
}
