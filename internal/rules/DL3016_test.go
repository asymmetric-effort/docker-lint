// file: internal/rules/DL3016_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/asymmetric-effort/docker-lint/internal/ir"
	"reflect"
)

// TestIntegrationPinNpmVersionID validates rule identity.
func TestIntegrationPinNpmVersionID(t *testing.T) {
	if NewPinNpmVersion().ID() != "DL3016" {
		t.Fatalf("unexpected id")
	}
}

// TestIntegrationPinNpmVersionViolation detects unpinned packages.
func TestIntegrationPinNpmVersionViolation(t *testing.T) {
	src := "FROM alpine\nRUN npm install express\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewPinNpmVersion()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
}

// TestIntegrationPinNpmVersionClean ensures compliant installs pass.
func TestIntegrationPinNpmVersionClean(t *testing.T) {
	src := "FROM alpine\nRUN npm install express@4.18.0\n"
	res, err := parser.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc, err := ir.BuildDocument("Dockerfile", res.AST)
	if err != nil {
		t.Fatalf("build document: %v", err)
	}
	r := NewPinNpmVersion()
	findings, err := r.Check(context.Background(), doc)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

// TestIntegrationPinNpmVersionNilDocument ensures graceful handling of nil input.
func TestIntegrationPinNpmVersionNilDocument(t *testing.T) {
	r := NewPinNpmVersion()
	if findings, err := r.Check(context.Background(), nil); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on nil doc: %v %v", findings, err)
	}
	if findings, err := r.Check(context.Background(), &ir.Document{}); err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings on empty doc: %v %v", findings, err)
	}
}

// TestUnitTrimFlag verifies flag normalization.
func TestUnitTrimFlag(t *testing.T) {
	cases := map[string]string{
		"--loglevel":      "loglevel",
		"--loglevel=info": "loglevel",
		"-v":              "v",
	}
	for in, want := range cases {
		if got := trimFlag(in); got != want {
			t.Fatalf("trimFlag(%q)=%q, want %q", in, got, want)
		}
	}
}

// TestUnitNpmInstallPackages validates package extraction.
func TestUnitNpmInstallPackages(t *testing.T) {
	tokens := []string{"npm", "--loglevel", "warn", "install", "--loglevel", "info", "left", "--loglevel=warn", "right"}
	pkgs := npmInstallPackages(tokens)
	if !reflect.DeepEqual(pkgs, []string{"left", "right"}) {
		t.Fatalf("unexpected packages: %#v", pkgs)
	}
}

// TestUnitNpmVersionFixed exercises version detection logic.
func TestUnitNpmVersionFixed(t *testing.T) {
	cases := []struct {
		pkg  string
		want bool
	}{
		{"express@1.0.0", true},
		{"@scope/pkg@2.3.4", true},
		{"git://example.com/repo.git#v1.0.0", true},
		{"git://example.com/repo.git", false},
		{"https://example.com/pkg.tgz", true},
		{"./localdir", true},
		{"express", false},
	}
	for _, c := range cases {
		if got := npmVersionFixed(c.pkg); got != c.want {
			t.Fatalf("npmVersionFixed(%q)=%v, want %v", c.pkg, got, c.want)
		}
	}
}
