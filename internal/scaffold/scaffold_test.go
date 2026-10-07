package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kuyacarlo/gitforge-toolkit/internal/config"
)

func TestRenderProject(t *testing.T) {
	stacks := []config.StackType{
		config.StackGo,
		config.StackPython,
		config.StackRust,
		config.StackC,
		config.StackJava,
	}

	for _, stack := range stacks {
		t.Run(string(stack), func(t *testing.T) {
			tmpDir := t.TempDir()
			cfg := &config.ProjectConfig{
				Name:      "sample-" + string(stack),
				Language:  stack,
				PreCommit: true,
				Workflows: true,
				Container: true,
			}

			if err := RenderProject(tmpDir, cfg); err != nil {
				t.Fatalf("RenderProject failed: %v", err)
			}

			// 1. Verify pre-commit
			pcBytes, err := os.ReadFile(filepath.Join(tmpDir, ".pre-commit-config.yaml"))
			if err != nil {
				t.Fatalf("missing .pre-commit-config.yaml: %v", err)
			}
			if len(pcBytes) == 0 {
				t.Error("empty pre-commit file")
			}

			// 2. Verify workflow
			wfBytes, err := os.ReadFile(filepath.Join(tmpDir, ".github", "workflows", "ci.yml"))
			if err != nil {
				t.Fatalf("missing ci.yml: %v", err)
			}
			if !strings.Contains(string(wfBytes), "sample-"+string(stack)) {
				t.Errorf("workflow missing project name: %s", string(wfBytes))
			}

			// 3. Verify Dockerfile
			dfBytes, err := os.ReadFile(filepath.Join(tmpDir, "Dockerfile"))
			if err != nil {
				t.Fatalf("missing Dockerfile: %v", err)
			}
			if len(dfBytes) == 0 {
				t.Error("empty Dockerfile")
			}

			// 4. Verify config saved
			if _, err := os.ReadFile(filepath.Join(tmpDir, "gitforge.yaml")); err != nil {
				t.Fatalf("missing gitforge.yaml: %v", err)
			}
		})
	}
}
