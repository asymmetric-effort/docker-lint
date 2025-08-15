// file: cmd/docker-lint/main.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

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
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}(f)

	res, err := parser.Parse(f)
	if err != nil {
		return err
	}
	doc, err := ir.BuildDocument(args[0], res.AST)
	if err != nil {
		return err
	}

	reg := engine.NewRegistry()
	reg.Register(rules.NewNoLatestTag())

	findings, err := reg.Run(context.Background(), doc)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(findings)
}
