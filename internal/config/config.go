// file: internal/config/config.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Config represents docker-lint configuration settings.
//
// The structure matches the configuration options used by hadolint so that
// users can reuse existing `.hadolint.yaml` files. Only a subset of hadolint's
// options are currently consumed by docker-lint.
type Config struct {
	// Ignored lists rule IDs that should be skipped globally during linting.
	Ignored []string `yaml:"ignored"`

	// Override remaps rule IDs to a severity level, keyed by that level.
	Override map[string][]string `yaml:"override"`

	// FailureThreshold defines the minimum severity that causes a failure.
	FailureThreshold string `yaml:"failure-threshold"`

	// TrustedRegistries specifies registries considered secure for FROM instructions.
	TrustedRegistries []string `yaml:"trustedRegistries"`

	// StrictLabels toggles enforcement of a configured label schema.
	StrictLabels bool `yaml:"strict-labels"`

	// LabelSchema maps required label keys to a description or type.
	LabelSchema map[string]string `yaml:"label-schema"`
}

// Load reads the configuration from the given YAML file path.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// IsIgnored reports whether the given rule ID is globally ignored.
func (c *Config) IsIgnored(rule string) bool {
	if c == nil {
		return false
	}
	for _, r := range c.Ignored {
		if r == rule {
			return true
		}
	}
	return false
}
