package main

import (
	"github.com/sardonyx001/unlinked/pkg/types"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [urls...]",
	Short: "Check specific URLs for broken links",
	Long: `Check one or more specific URLs to see if they are accessible.
This command validates individual URLs without crawling for additional links.

Examples:
  # Check a single URL
  unlinked check https://example.com

  # Check multiple URLs
  unlinked check https://example.com https://golang.org

  # Check URLs from stdin
  cat urls.txt | unlinked check --stdin

  # Save results as markdown
  unlinked check --output-format=markdown --output-file=report.md https://example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Force single mode
		cfg.Set("mode", types.ModeSingle)
		return runCheck(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Output flags
	checkCmd.Flags().StringVarP(&flagOutputFormat, "output-format", "f", "plaintext", "output format: plaintext, markdown, html, json")
	checkCmd.Flags().StringVarP(&flagOutputFile, "output-file", "o", "", "output file (default is stdout)")

	// Behavior flags
	checkCmd.Flags().IntVarP(&flagConcurrency, "concurrency", "c", 10, "number of concurrent checks")
	checkCmd.Flags().IntVarP(&flagTimeout, "timeout", "t", 30, "timeout in seconds for each request")
	checkCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	checkCmd.Flags().BoolVar(&flagNoProgress, "no-progress", false, "disable progress display")
	checkCmd.Flags().BoolVar(&flagStdin, "stdin", false, "read URLs from stdin")
}
