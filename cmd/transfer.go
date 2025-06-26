package cmd

import (
	"fmt"

	"github.com/lilmonk/elasticdump/internal/transfer"
	"github.com/spf13/cobra"
)

var (
	input       string
	output      string
	dataType    string
	limit       int
	concurrency int
	format      string
	scrollSize  int
	username    string
	password    string
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer data between Elasticsearch clusters",
	Long: `Transfer data, mappings, or settings between Elasticsearch clusters.
This command supports various transfer types and can handle large datasets efficiently.`,
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
	rootCmd.AddCommand(transferCmd)

	// Transfer flags
	transferCmd.Flags().StringVarP(&input, "input", "i", "", "Source Elasticsearch cluster or index (required)")
	transferCmd.Flags().StringVarP(&output, "output", "o", "", "Destination Elasticsearch cluster or index (required)")
	transferCmd.Flags().StringVarP(&dataType, "type", "t", "data", "Type of data to transfer (data, mapping, settings)")
	transferCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number of records to transfer (0 = no limit)")
	transferCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 4, "Number of concurrent operations")
	transferCmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, ndjson)")
	transferCmd.Flags().IntVarP(&scrollSize, "scrollSize", "s", 1000, "Size of the scroll for large datasets")
	transferCmd.Flags().StringVarP(&username, "username", "u", "", "Elasticsearch username (optional)")
	transferCmd.Flags().StringVarP(&password, "password", "p", "", "Elasticsearch password (optional)")

	// Mark required flags
	transferCmd.MarkFlagRequired("input")
	transferCmd.MarkFlagRequired("output")
}
