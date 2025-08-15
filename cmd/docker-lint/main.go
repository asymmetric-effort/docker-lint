// file: cmd/docker-lint/main.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/sam-caldwell/ansi"

	"github.com/asymmetric-effort/docker-lint/internal/config"
	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"github.com/asymmetric-effort/docker-lint/internal/rules"
)

// usageText describes the command line usage for the application.
const usageText = "usage: docker-lint [-c file] <Dockerfile>"

// printUsage writes the CLI usage information to the provided writer.
func printUsage(out io.Writer) {
	fmt.Fprintln(out, usageText)
}

// main is the CLI entry point.
func main() {
	color := true
	var args []string
	for _, a := range os.Args[1:] {
		if a == "--no-color" {
			color = false
			continue
		}
		args = append(args, a)
	}
	if err := run(args, os.Stdout, os.Stderr, color); err != nil {
		printError(os.Stderr, color, err)
		os.Exit(1)
	}
}

// run executes the linter for the provided arguments and writes JSON findings.
//
// In addition to the JSON output, run emits a human-readable summary to errOut.
// When color is true, the summary uses ANSI colors.
func run(args []string, out io.Writer, errOut io.Writer, color bool) error {
	var (
		files      []string
		configPath string
	)
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-h", "--help":
			printUsage(out)
			return nil
		case "-c", "--config":
			if i+1 >= len(args) {
				return fmt.Errorf("missing config file after %s", a)
			}
			configPath = args[i+1]
			i++
		default:
			files = append(files, a)
		}
	}
	if len(files) == 0 {
		return fmt.Errorf(usageText)
	}

	var cfg config.Config
	if configPath != "" {
		c, err := config.Load(configPath)
		if err != nil {
			return err
		}
		cfg = *c
	} else {
		if _, err := os.Stat(".docker-lint.yaml"); err == nil {
			c, err := config.Load(".docker-lint.yaml")
			if err != nil {
				return err
			}
			cfg = *c
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	files, err := expandPaths(files)
	if err != nil {
		return err
	}

	reg := engine.NewRegistry()
	registerRules(reg, cfg.Exclusions)

	ctx := context.Background()
	var all []engine.Finding
	for _, path := range files {
		fnds, err := lintFile(ctx, reg, path)
		if err != nil {
			return err
		}
		all = append(all, fnds...)
	}
	if err := json.NewEncoder(out).Encode(all); err != nil {
		return err
	}
	printFindings(errOut, all, color)
	return nil
}

// registerRules adds built-in rules to reg, skipping any whose IDs appear in exclusions.
func registerRules(reg *engine.Registry, exclusions []string) {
	skip := map[string]struct{}{}
	for _, id := range exclusions {
		skip[id] = struct{}{}
	}
	all := []engine.Rule{
		rules.NewNoLatestTag(),
		rules.NewRequireOSVersionTag(),
		rules.NewAptNoInstallRecommends(),
		rules.NewAptPin(),
		rules.NewAptListsCleanup(),
		rules.NewDnfNoUpgrade(),
		rules.NewDnfCacheCleanup(),
	}
	for _, r := range all {
		if _, ok := skip[r.ID()]; ok {
			continue
		}
		reg.Register(r)
	}
}

// printFindings writes a human-readable summary of findings to errOut.
func printFindings(errOut io.Writer, fnds []engine.Finding, color bool) {
	if len(fnds) == 0 {
		if color {
			fmt.Fprintf(errOut, "%sNo issues found%s\n", ansi.CodeFgGreen, ansi.CodeReset)
		} else {
			fmt.Fprintln(errOut, "No issues found")
		}
		return
	}
	for _, f := range fnds {
		if color {
			fmt.Fprintf(errOut, "%s%s (rule: %s, line: %d)%s\n", ansi.CodeFgYellow, f.Message, f.RuleID, f.Line, ansi.CodeReset)
		} else {
			fmt.Fprintf(errOut, "%s (rule: %s, line: %d)\n", f.Message, f.RuleID, f.Line)
		}
	}
}

// printError writes an error message to errOut, optionally colorized.
func printError(errOut io.Writer, color bool, err error) {
	if color {
		fmt.Fprintf(errOut, "%s%v%s\n", ansi.CodeFgRed, err, ansi.CodeReset)
	} else {
		fmt.Fprintln(errOut, err)
	}
}

// expandPaths resolves glob patterns into file paths.
func expandPaths(patterns []string) ([]string, error) {
	var files []string
	for _, p := range patterns {
		matches, err := doublestar.FilepathGlob(p)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			files = append(files, p)
			continue
		}
		files = append(files, matches...)
	}
	return files, nil
}

// lintFile lints a single Dockerfile and returns any findings.
func lintFile(ctx context.Context, reg *engine.Registry, path string) (fnds []engine.Finding, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	res, err := parser.Parse(f)
	if err != nil {
		return nil, err
	}
	doc, err := ir.BuildDocument(path, res.AST)
	if err != nil {
		return nil, err
	}
	return reg.Run(ctx, doc)
}
