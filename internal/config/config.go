// file: internal/config/config.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents docker-lint configuration.
//
// Config currently supports file-based rule exclusions.
type Config struct {
	Exclude map[string][]string `yaml:"exclude"`
}

// Load reads the configuration from the provided path.
//
// Load returns an error if the file cannot be read or parsed.
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

// IsRuleExcluded reports whether the given rule is excluded for the file.
func (c *Config) IsRuleExcluded(path, rule string) bool {
	if c == nil {
		return false
	}
	base := filepath.Base(path)
	rules, ok := c.Exclude[base]
	if !ok {
		return false
	}
	for _, r := range rules {
		if r == rule {
			return true
		}
	}
	return false
}
