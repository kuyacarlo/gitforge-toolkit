package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/kuyacarlo/gitforge-toolkit/internal/checker"
	"github.com/kuyacarlo/gitforge-toolkit/internal/config"
	"github.com/spf13/cobra"
)

var checkConfigPath string

var checkCmd = &cobra.Command{
	Use:   "check [directory]",
	Short: "Check local toolchain and binary dependencies",
	Long:  `Verify required compilers, runtimes, and tools defined in gitforge.yaml are installed and meet version requirements.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := checkConfigPath
		if targetPath == "" {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			targetPath = filepath.Join(dir, "gitforge.yaml")
		}

		cfg, err := config.LoadConfig(targetPath)
		if err != nil {
			return fmt.Errorf("could not load configuration from %s: %w", targetPath, err)
		}

		fmt.Printf("Auditing toolchain requirements for '%s' (%s)...\n", cfg.Name, cfg.Language)

		if len(cfg.Tools) == 0 {
			fmt.Println("No tool requirements defined in configuration.")
			return nil
		}

		results := checker.CheckProject(cfg)
		allPassed := true

		fmt.Println("\nTOOLCHAIN DIAGNOSTIC MATRIX:")
		fmt.Printf("%-15s %-12s %-15s %-15s %s\n", "BINARY", "STATUS", "DETECTED", "REQUIRED", "DETAILS")
		fmt.Println("--------------------------------------------------------------------------------")

		for _, r := range results {
			statusStr := "PASSED"
			if !r.Passed {
				statusStr = "FAILED"
				allPassed = false
			}

			installedStr := r.Version
			if !r.Installed {
				installedStr = "NOT FOUND"
			}

			reqStr := r.Required
			if reqStr == "" {
				reqStr = "(any)"
			}

			details := ""
			if r.Error != "" {
				details = r.Error
			}

			fmt.Printf("%-15s %-12s %-15s %-15s %s\n", r.Binary, statusStr, installedStr, reqStr, details)
		}

		fmt.Println()
		if !allPassed {
			return fmt.Errorf("toolchain check failed: some dependencies are missing or incompatible")
		}

		fmt.Println("All toolchain requirements are satisfied!")
		return nil
	},
}

func init() {
	checkCmd.Flags().StringVarP(&checkConfigPath, "config", "c", "", "Path to gitforge.yaml config file")
	rootCmd.AddCommand(checkCmd)
}
