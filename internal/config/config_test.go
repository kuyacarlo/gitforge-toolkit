package config

import (
	"path/filepath"
	"testing"
)

func TestLoadAndSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "gitforge.yaml")

	original := &ProjectConfig{
		Name:      "test-service",
		Language:  StackGo,
		Version:   "1.22",
		PreCommit: true,
		Workflows: true,
		Container: true,
		Tools: []ToolRequirement{
			{Binary: "go", MinVersion: "1.22.0"},
			{Binary: "podman"},
		},
	}

	if err := original.Save(cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Name != original.Name {
		t.Errorf("expected name %s, got %s", original.Name, loaded.Name)
	}
	if loaded.Language != original.Language {
		t.Errorf("expected language %s, got %s", original.Language, loaded.Language)
	}
	if len(loaded.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(loaded.Tools))
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/gitforge.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
