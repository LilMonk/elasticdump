package cmd

import (
	"fmt"

	"github.com/lilmonk/elasticdump/internal/transfer"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup Elasticsearch data to file",
	Long: `Backup Elasticsearch data, mappings, or settings to a file.
This is a convenient wrapper around the transfer command for backup operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if input == "" {
			return fmt.Errorf("input is required")
		}
		if output == "" {
			return fmt.Errorf("output is required")
		}

		config := transfer.Config{
			Input:       input,
			Output:      output,
			Type:        dataType,
			Limit:       limit,
			Concurrency: concurrency,
			Format:      format,
			ScrollSize:  scrollSize,
			Verbose:     verbose,
			Username:    username,
			Password:    password,
		}

		return transfer.Run(config)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Backup flags (reuse the same variables from transfer command)
	backupCmd.Flags().StringVarP(&input, "input", "i", "", "Source Elasticsearch cluster or index (required)")
	backupCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (required)")
	backupCmd.Flags().StringVarP(&dataType, "type", "t", "data", "Type of data to backup (data, mapping, settings)")
	backupCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number of records to backup (0 = no limit)")
	backupCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 4, "Number of concurrent operations")
	backupCmd.Flags().StringVarP(&format, "format", "f", "ndjson", "Output format (json, ndjson)")
	backupCmd.Flags().IntVarP(&scrollSize, "scrollSize", "s", 1000, "Size of the scroll for large datasets")
	backupCmd.Flags().StringVarP(&username, "username", "u", "", "Elasticsearch username (optional)")
	backupCmd.Flags().StringVarP(&password, "password", "p", "", "Elasticsearch password (optional)")

	// Mark required flags
	backupCmd.MarkFlagRequired("input")
	backupCmd.MarkFlagRequired("output")
}
