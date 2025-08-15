// file: internal/rules/DL3042_helpers_test.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package rules

import "testing"

// TestPackageManagerFamily verifies detection of package manager families.
func TestPackageManagerFamily(t *testing.T) {
	cases := map[string]struct {
		seg    []string
		expect string
	}{
		"apt-get":   {[]string{"apt-get", "update"}, "apt"},
		"apt":       {[]string{"apt", "install"}, "apt"},
		"apk":       {[]string{"apk", "add"}, "apk"},
		"dnf":       {[]string{"dnf", "upgrade"}, "dnf"},
		"microdnf":  {[]string{"microdnf", "clean"}, "dnf"},
		"yum":       {[]string{"yum", "remove"}, "yum"},
		"zypper":    {[]string{"zypper", "install"}, "zypper"},
		"flag-skip": {[]string{"apt-get", "-y", "install"}, "apt"},
		"unknown":   {[]string{"echo", "hi"}, ""},
		"short":     {[]string{"apt-get"}, ""},
		"missing":   {[]string{"apt-get", "-y"}, ""},
	}
	for name, tc := range cases {
		if got := packageManagerFamily(tc.seg); got != tc.expect {
			t.Fatalf("%s: expected %s got %s", name, tc.expect, got)
		}
	}
}
