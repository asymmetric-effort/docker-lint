// Copyright (c) 2025 Asymmetric Effort, LLC

package packaging

import "testing"

// TestName verifies that package file names are generated correctly.
func TestName(t *testing.T) {
	cases := []struct {
		version, goos, arch, ext, expect string
	}{
		{"v1.0.0", "linux", "amd64", "deb", "docker-lint_v1.0.0_linux_amd64.deb"},
		{"v1.0.0", "linux", "arm64", "rpm", "docker-lint_v1.0.0_linux_arm64.rpm"},
		{"v1.0.0", "windows", "amd64", "nupkg", "docker-lint_v1.0.0_windows_amd64.nupkg"},
	}

	for _, c := range cases {
		name := Name(c.version, c.goos, c.arch, c.ext)
		if name != c.expect {
			t.Fatalf("expected %s, got %s", c.expect, name)
		}
	}
}
