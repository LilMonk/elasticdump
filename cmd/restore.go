package cmd

import (
	"fmt"

	"github.com/lilmonk/elasticdump/internal/restore"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Elasticsearch data from file",
	Long: `Restore Elasticsearch data, mappings, or settings from a backup file.
This command reads data from a file and imports it into an Elasticsearch cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if input == "" {
			return fmt.Errorf("input file is required")
		}
		if output == "" {
			return fmt.Errorf("output cluster is required")
		}

		config := restore.Config{
			Input:       input,
			Output:      output,
			Type:        dataType,
			Concurrency: concurrency,
			Verbose:     verbose,
			Username:    username,
			Password:    password,
		}

		return restore.Run(config)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Restore flags
	restoreCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path (required)")
	restoreCmd.Flags().StringVarP(&output, "output", "o", "", "Destination Elasticsearch cluster or index (required)")
	restoreCmd.Flags().StringVarP(&dataType, "type", "t", "data", "Type of data to restore (data, mapping, settings)")
	restoreCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 4, "Number of concurrent operations")
	restoreCmd.Flags().StringVarP(&username, "username", "u", "", "Elasticsearch username (optional)")
	restoreCmd.Flags().StringVarP(&password, "password", "p", "", "Elasticsearch password (optional)")

	// Mark required flags
	restoreCmd.MarkFlagRequired("input")
	restoreCmd.MarkFlagRequired("output")
}
