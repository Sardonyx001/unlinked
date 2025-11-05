package main

import (
	"github.com/sardonyx001/unlinked/pkg/types"
	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl [url]",
	Short: "Crawl a website and check all discovered links",
	Long: `Crawl a website starting from the given URL and check all discovered links.
This command recursively follows links up to the specified depth and validates each one.

Examples:
  # Crawl a website with default settings
  unlinked crawl https://example.com

  # Crawl with custom depth
  unlinked crawl --max-depth=2 https://example.com

  # Fast crawl with high concurrency
  unlinked crawl --concurrency=50 --max-depth=3 https://example.com

  # Crawl and save HTML report
  unlinked crawl --output-format=html --output-file=report.html https://example.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Force crawler mode
		cfg.Set("mode", types.ModeCrawler)
		return runCheck(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	// Crawler-specific flags
	crawlCmd.Flags().IntVar(&flagMaxDepth, "max-depth", 3, "maximum crawl depth")

	// Output flags
	crawlCmd.Flags().StringVarP(&flagOutputFormat, "output-format", "f", "plaintext", "output format: plaintext, markdown, html, json")
	crawlCmd.Flags().StringVarP(&flagOutputFile, "output-file", "o", "", "output file (default is stdout)")

	// Behavior flags
	crawlCmd.Flags().IntVarP(&flagConcurrency, "concurrency", "c", 10, "number of concurrent checks")
	crawlCmd.Flags().IntVarP(&flagTimeout, "timeout", "t", 30, "timeout in seconds for each request")
	crawlCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	crawlCmd.Flags().BoolVar(&flagNoProgress, "no-progress", false, "disable progress display")
}
