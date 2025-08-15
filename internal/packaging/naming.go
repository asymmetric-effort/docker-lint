// Copyright (c) 2025 Asymmetric Effort, LLC

package packaging

// Name generates a package file name for the given parameters.
// The format is "docker-lint_<version>_<os>_<arch>.<ext>".
func Name(version, goos, arch, ext string) string {
	return "docker-lint_" + version + "_" + goos + "_" + arch + "." + ext
}
