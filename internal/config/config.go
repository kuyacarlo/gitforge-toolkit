package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type StackType string

const (
	StackGo     StackType = "go"
	StackPython StackType = "python"
	StackRust   StackType = "rust"
	StackC      StackType = "c"
	StackJava   StackType = "java"
)

type ToolRequirement struct {
	Binary     string `yaml:"binary"`
	MinVersion string `yaml:"min_version,omitempty"`
}

type ProjectConfig struct {
	Name        string            `yaml:"name"`
	Language    StackType         `yaml:"language"`
	Version     string            `yaml:"version,omitempty"`
	Tools       []ToolRequirement `yaml:"tools,omitempty"`
	PreCommit   bool              `yaml:"pre_commit"`
	Workflows   bool              `yaml:"workflows"`
	Container   bool              `yaml:"container"`
}

func LoadConfig(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing yaml config: %w", err)
	}

	return &cfg, nil
}

func (c *ProjectConfig) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling yaml config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config file %s: %w", path, err)
	}
	return nil
}
