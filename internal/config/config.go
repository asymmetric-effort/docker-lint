// file: internal/config/config.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents docker-lint configuration settings.
//
// Exclusions lists rule IDs that should be skipped during linting.
type Config struct {
	Exclusions []string `yaml:"exclusions"`
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
