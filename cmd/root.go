package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitforge",
	Short: "gitforge-toolkit is a developer tooling and scaffolding CLI",
	Long:  `gitforge is a comprehensive CLI toolkit to initialize project scaffolding and check local toolchains.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here if needed
}
