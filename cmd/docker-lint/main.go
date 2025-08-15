// file: cmd/docker-lint/main.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/engine"
	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"github.com/asymmetric-effort/docker-lint/internal/rules"
)

// usageText describes the command line usage for the application.
const usageText = "usage: docker-lint <Dockerfile>"

// printUsage writes the CLI usage information to the provided writer.
func printUsage(out io.Writer) {
	fmt.Fprintln(out, usageText)
}

// main is the CLI entry point.
func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run executes the linter for the provided arguments and writes JSON findings.
func run(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf(usageText)
	}
	if args[0] == "-h" || args[0] == "--help" {
		printUsage(out)
		return nil
	}

	files, err := expandPaths(args)
	if err != nil {
		return err
	}

	reg := engine.NewRegistry()
	reg.Register(rules.NewNoLatestTag())
	reg.Register(rules.NewAptListsCleanup())

	ctx := context.Background()
	var all []engine.Finding
	for _, path := range files {
		fnds, err := lintFile(ctx, reg, path)
		if err != nil {
			return err
		}
		all = append(all, fnds...)
	}
	return json.NewEncoder(out).Encode(all)
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
