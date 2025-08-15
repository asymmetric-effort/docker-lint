// file: internal/version/version.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
// Package version provides application build version information.
package version

// Current holds the current application version.
// This value may be overridden at build time using -ldflags.
var Current = "dev"
