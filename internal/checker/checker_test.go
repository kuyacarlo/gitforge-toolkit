package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kuyacarlo/gitforge-toolkit/internal/config"
)

func TestParseVer(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"1.22.0", []int{1, 22, 0}},
		{"v1.22.4", []int{1, 22, 4}},
		{"3.12", []int{3, 12}},
		{"21", []int{21}},
	}

	for _, tt := range tests {
		res := parseVer(tt.input)
		if len(res) != len(tt.expected) {
			t.Errorf("parseVer(%s) returned length %d, expected %d", tt.input, len(res), len(tt.expected))
			continue
		}
		for i := range res {
			if res[i] != tt.expected[i] {
				t.Errorf("parseVer(%s)[%d] = %d, expected %d", tt.input, i, res[i], tt.expected[i])
			}
		}
	}
}

func TestCompareSemVer(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.22.0", "1.22.0", 0},
		{"1.22.1", "1.22.0", 1},
		{"1.21.9", "1.22.0", -1},
		{"1.22", "1.22.0", 0},
		{"2.0.0", "1.99.99", 1},
		{"v1.20", "1.19.5", 1},
	}

	for _, tt := range tests {
		got := compareSemVer(tt.v1, tt.v2)
		if got != tt.expected {
			t.Errorf("compareSemVer(%s, %s) = %d, expected %d", tt.v1, tt.v2, got, tt.expected)
		}
	}
}

func TestCheckToolMissing(t *testing.T) {
	res := CheckTool(config.ToolRequirement{
		Binary:     "nonexistent-tool-binary-xyz",
		MinVersion: "1.0.0",
	})

	if res.Installed {
		t.Errorf("expected tool to be not installed")
	}
	if res.Passed {
		t.Errorf("expected missing tool to fail check")
	}
	if res.Error == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestCheckToolMockBinary(t *testing.T) {
	tmpDir := t.TempDir()
	mockBin := filepath.Join(tmpDir, "mock-tool")

	script := "#!/bin/sh\necho \"mock-tool version 2.4.1\"\n"
	if err := os.WriteFile(mockBin, []byte(script), 0755); err != nil {
		t.Fatalf("failed to create mock tool: %v", err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+origPath)
	defer os.Setenv("PATH", origPath)

	// Test 1: Passing version requirement
	res := CheckTool(config.ToolRequirement{
		Binary:     "mock-tool",
		MinVersion: "2.0.0",
	})
	if !res.Installed || !res.Passed || res.Version != "2.4.1" {
		t.Errorf("unexpected check result: %+v", res)
	}

	// Test 2: Failing version requirement
	resFail := CheckTool(config.ToolRequirement{
		Binary:     "mock-tool",
		MinVersion: "3.0.0",
	})
	if !resFail.Installed || resFail.Passed {
		t.Errorf("expected tool to fail version check: %+v", resFail)
	}

	// Test 3: No version requirement
	resNoMin := CheckTool(config.ToolRequirement{
		Binary: "mock-tool",
	})
	if !resNoMin.Installed || !resNoMin.Passed {
		t.Errorf("expected tool to pass when no min version is specified: %+v", resNoMin)
	}
}

func TestCheckProject(t *testing.T) {
	cfg := &config.ProjectConfig{
		Tools: []config.ToolRequirement{
			{Binary: "nonexistent-tool-1"},
			{Binary: "nonexistent-tool-2"},
		},
	}
	results := CheckProject(cfg)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
