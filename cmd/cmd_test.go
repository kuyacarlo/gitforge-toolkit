package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmd(t *testing.T) {
	tmpDir := t.TempDir()

	// Reset flags and command
	rootCmd.SetArgs([]string{"init", tmpDir, "--name", "test-project", "--lang", "go"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected init command to succeed, got %v", err)
	}

	// Verify generated files
	if _, err := os.Stat(filepath.Join(tmpDir, "gitforge.yaml")); os.IsNotExist(err) {
		t.Errorf("gitforge.yaml was not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".pre-commit-config.yaml")); os.IsNotExist(err) {
		t.Errorf(".pre-commit-config.yaml was not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".github", "workflows", "ci.yml")); os.IsNotExist(err) {
		t.Errorf(".github/workflows/ci.yml was not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "Dockerfile")); os.IsNotExist(err) {
		t.Errorf("Dockerfile was not created")
	}
}

func TestCheckCmd(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Initialize Go project in tmpDir
	rootCmd.SetArgs([]string{"init", tmpDir, "--name", "check-test", "--lang", "go"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// 2. Run check command on tmpDir (go should be present in our env)
	rootCmd.SetArgs([]string{"check", tmpDir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("check failed: %v", err)
	}
}
