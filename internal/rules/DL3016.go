package rules

/*
 * file: internal/rules/DL3016.go
 * (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
 */

import (
	"context"
	"strings"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
)

// pinNpmVersion enforces version pinning for npm install commands.
type pinNpmVersion struct{}

// NewPinNpmVersion constructs the rule.
func NewPinNpmVersion() engine.Rule { return pinNpmVersion{} }

// ID returns the rule identifier.
func (pinNpmVersion) ID() string { return "DL3016" }

// Check scans RUN instructions for unpinned npm install usage.
func (pinNpmVersion) Check(ctx context.Context, d *ir.Document) ([]engine.Finding, error) {
	var findings []engine.Finding
	if d == nil || d.AST == nil {
		return findings, nil
	}
	for _, n := range d.AST.Children {
		if !strings.EqualFold(n.Value, "run") {
			continue
		}
		segments := splitRunSegmentsNpm(n)
		for _, seg := range segments {
			if len(seg) == 0 {
				continue
			}
			if strings.ToLower(seg[0]) != "npm" {
				continue
			}
			if packages := npmInstallPackages(seg); len(packages) > 0 && !allVersionFixed(packages) {
				findings = append(findings, engine.Finding{
					RuleID:  "DL3016",
					Message: "Pin versions in npm. Instead of `npm install <package>` use `npm install <package>@<version>`",
					Line:    n.StartLine,
				})
			}
		}
	}
	return findings, nil
}

// npmInstallPackages returns package arguments to npm install if present.
// Flags like --loglevel are ignored.
func npmInstallPackages(tokens []string) []string {
	ignore := map[string]struct{}{"loglevel": {}}
	i := 1
	// skip flags before subcommand
	for i < len(tokens) {
		if !strings.HasPrefix(tokens[i], "-") {
			break
		}
		name := trimFlag(tokens[i])
		if _, ok := ignore[name]; ok && !strings.Contains(tokens[i], "=") {
			i += 2
		} else {
			i++
		}
	}
	if i >= len(tokens) || strings.ToLower(tokens[i]) != "install" {
		return nil
	}
	i++
	var packages []string
	for i < len(tokens) {
		tok := tokens[i]
		if strings.HasPrefix(tok, "-") {
			name := trimFlag(tok)
			if _, ok := ignore[name]; ok && !strings.Contains(tok, "=") {
				i += 2
				continue
			}
			i++
			continue
		}
		packages = append(packages, tok)
		i++
	}
	return packages
}

// splitRunSegmentsNpm tokenizes a RUN node and splits it into command segments.
func splitRunSegmentsNpm(n *parser.Node) [][]string {
	if n == nil || n.Next == nil {
		return nil
	}
	var tokens []string
	if n.Attributes != nil && n.Attributes["json"] {
		for tok := n.Next; tok != nil; tok = tok.Next {
			tokens = append(tokens, tok.Value)
		}
	} else {
		t, err := shlex.Split(n.Next.Value)
		if err != nil {
			return nil
		}
		tokens = t
	}
	var segments [][]string
	var current []string
	for _, tok := range tokens {
		switch tok {
		case "&&", "||", "|", ";":
			if len(current) > 0 {
				segments = append(segments, current)
				current = nil
			}
		default:
			current = append(current, tok)
		}
	}
	if len(current) > 0 {
		segments = append(segments, current)
	}
	return segments
}

// allVersionFixed returns true if all packages specify a version, tag, or commit.
func allVersionFixed(pkgs []string) bool {
	for _, p := range pkgs {
		if !npmVersionFixed(p) {
			return false
		}
	}
	return true
}

// npmVersionFixed determines if a package argument is pinned to a version.
func npmVersionFixed(pkg string) bool {
	if hasGitPrefix(pkg) {
		return isVersionedGit(pkg)
	}
	if hasTarballSuffix(pkg) {
		return true
	}
	if isFolder(pkg) {
		return true
	}
	return hasVersionSymbol(pkg)
}

func hasGitPrefix(pkg string) bool {
	prefixes := []string{"git://", "git+ssh://", "git+http://", "git+https://"}
	for _, p := range prefixes {
		if strings.HasPrefix(pkg, p) {
			return true
		}
	}
	return false
}

func hasTarballSuffix(pkg string) bool {
	suffixes := []string{".tar", ".tar.gz", ".tgz"}
	for _, s := range suffixes {
		if strings.HasSuffix(pkg, s) {
			return true
		}
	}
	return false
}

func isFolder(pkg string) bool {
	prefixes := []string{"/", "./", "../", "~/"}
	for _, p := range prefixes {
		if strings.HasPrefix(pkg, p) {
			return true
		}
	}
	return false
}

func isVersionedGit(pkg string) bool { return strings.Contains(pkg, "#") }

func hasVersionSymbol(pkg string) bool {
	if strings.HasPrefix(pkg, "@") {
		if idx := strings.Index(pkg, "/"); idx != -1 {
			pkg = pkg[idx+1:]
		}
	}
	return strings.Contains(pkg, "@")
}

// trimFlag normalizes a flag token for comparison.
func trimFlag(flag string) string {
	flag = strings.TrimLeft(flag, "-")
	if idx := strings.Index(flag, "="); idx != -1 {
		flag = flag[:idx]
	}
	return flag
}
