package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuyacarlo/gitforge-toolkit/internal/config"
	"github.com/kuyacarlo/gitforge-toolkit/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	initName      string
	initLang      string
	initPreCommit bool
	initWorkflows bool
	initContainer bool
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Scaffold a new repository structure and configs",
	Long:  `Initialize a new repository with pre-commit hooks, CI workflows, Dockerfiles, and gitforge.yaml.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}

		name := initName
		if name == "" {
			abs, err := filepath.Abs(targetDir)
			if err == nil {
				name = filepath.Base(abs)
			} else {
				name = "my-project"
			}
		}

		stack := config.StackType(initLang)
		switch stack {
		case config.StackGo, config.StackPython, config.StackRust, config.StackC, config.StackJava:
		default:
			return fmt.Errorf("unsupported language/stack '%s'. Supported stacks: go, python, rust, c, java", initLang)
		}

		cfg := &config.ProjectConfig{
			Name:      name,
			Language:  stack,
			PreCommit: initPreCommit,
			Workflows: initWorkflows,
			Container: initContainer,
		}

		// Configure default tool requirements based on stack
		switch stack {
		case config.StackGo:
			cfg.Tools = []config.ToolRequirement{
				{Binary: "go", MinVersion: "1.20"},
			}
		case config.StackPython:
			cfg.Tools = []config.ToolRequirement{
				{Binary: "python3", MinVersion: "3.10"},
			}
		case config.StackRust:
			cfg.Tools = []config.ToolRequirement{
				{Binary: "cargo", MinVersion: "1.70"},
			}
		case config.StackC:
			cfg.Tools = []config.ToolRequirement{
				{Binary: "gcc"},
				{Binary: "make"},
			}
		case config.StackJava:
			cfg.Tools = []config.ToolRequirement{
				{Binary: "java", MinVersion: "17"},
			}
		}

		if err := scaffold.RenderProject(targetDir, cfg); err != nil {
			return fmt.Errorf("scaffolding project: %w", err)
		}

		fmt.Printf("Successfully initialized gitforge project '%s' (%s) in %s\n", name, initLang, targetDir)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initName, "name", "n", "", "Project name (default: directory name)")
	initCmd.Flags().StringVarP(&initLang, "lang", "l", "go", "Project language stack (go, python, rust, c, java)")
	initCmd.Flags().BoolVar(&initPreCommit, "pre-commit", true, "Generate .pre-commit-config.yaml")
	initCmd.Flags().BoolVar(&initWorkflows, "workflows", true, "Generate GitHub Actions CI workflow")
	initCmd.Flags().BoolVar(&initContainer, "container", true, "Generate Dockerfile")

	rootCmd.AddCommand(initCmd)
}
