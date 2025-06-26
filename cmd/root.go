package cmd

import (
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elasticdump",
	Short: "A Go CLI tool to migrate data between Elasticsearch clusters",
	Long: `Elasticdump is a powerful CLI tool built in Go that allows you to:
- Migrate data between Elasticsearch clusters
- Backup and restore data
- Handle large datasets efficiently
- Support multiple output formats (JSON, NDJSON, etc.)
- Perform multi-threaded operations for faster processing`,
	Version: "dev",
}

// SetVersion sets the version for the root command
func SetVersion(version string) {
	rootCmd.Version = version
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
