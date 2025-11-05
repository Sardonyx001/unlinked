package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sardonyx001/unlinked/internal/checker"
	"github.com/sardonyx001/unlinked/internal/config"
	"github.com/sardonyx001/unlinked/internal/output"
	"github.com/sardonyx001/unlinked/internal/ui"
	"github.com/sardonyx001/unlinked/pkg/types"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Manager

	// Flags
	flagMode           string
	flagOutputFormat   string
	flagOutputFile     string
	flagConcurrency    int
	flagTimeout        int
	flagMaxDepth       int
	flagVerbose        bool
	flagNoProgress     bool
	flagStdin          bool
)

var rootCmd = &cobra.Command{
	Use:   "unlinked",
	Short: "A modern dead link checker CLI",
	Long: `Unlinked is a fast, modern dead link checker that can check individual URLs
or crawl websites to find broken links. It supports multiple output formats
and provides a beautiful terminal UI.

Commands:
  check       Check specific URLs for broken links
  crawl       Crawl a website and check all discovered links
  version     Print version information
  completion  Generate shell completion scripts
  help        Show help for any command
  list        List all available commands

Examples:
  # Check specific URLs
  unlinked check https://example.com

  # Crawl a website
  unlinked crawl https://example.com

  # Show version
  unlinked version

  # Generate bash completion
  unlinked completion bash`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Config file flag (persistent across all subcommands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/unlinked/config.yaml)")
}

func initConfig() {
	cfg = config.New()
	if err := cfg.Load(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
}

// collectURLs collects URLs from arguments and/or stdin
func collectURLs(args []string) ([]string, error) {
	var urls []string

	// Read from stdin if flag is set
	if flagStdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				urls = append(urls, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading from stdin: %w", err)
		}
	}

	// Add URLs from arguments
	urls = append(urls, args...)

	return urls, nil
}

func Execute() error {
	return rootCmd.Execute()
}

// runCheck is the shared logic for check and crawl commands
func runCheck(cmd *cobra.Command, args []string) error {
	// Apply command-line flags to config
	applyFlags(cmd)

	// Collect URLs to check
	urls, err := collectURLs(args)
	if err != nil {
		return fmt.Errorf("failed to collect URLs: %w", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs provided. Use --stdin to read from stdin or provide URLs as arguments")
	}

	// Create checker
	c, err := checker.New(cfg.Get())
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

	// Set up UI
	var result *types.CheckResult
	if cfg.Get().ShowProgress {
		result, err = runWithUI(c, urls)
	} else {
		result, err = runWithoutUI(c, urls)
	}

	if err != nil {
		return fmt.Errorf("check failed: %w", err)
	}

	// Format and output results
	if err := outputResults(result); err != nil {
		return fmt.Errorf("failed to output results: %w", err)
	}

	// Exit with error code if issues found
	if result.TotalDead > 0 || result.TotalErrors > 0 {
		os.Exit(1)
	}

	return nil
}

func applyFlags(cmd *cobra.Command) {
	if cmd.Flags().Changed("output-format") {
		cfg.Set("output_format", types.OutputFormat(flagOutputFormat))
	}
	if cmd.Flags().Changed("output-file") {
		cfg.Set("output_file", flagOutputFile)
	}
	if cmd.Flags().Changed("concurrency") {
		cfg.Set("concurrency", flagConcurrency)
	}
	if cmd.Flags().Changed("timeout") {
		cfg.Set("timeout", flagTimeout)
	}
	if cmd.Flags().Changed("max-depth") {
		cfg.Set("max_depth", flagMaxDepth)
	}
	if cmd.Flags().Changed("verbose") {
		cfg.Set("verbose", flagVerbose)
	}
	if cmd.Flags().Changed("no-progress") {
		cfg.Set("show_progress", !flagNoProgress)
	}
}

func runWithUI(c *checker.Checker, urls []string) (*types.CheckResult, error) {
	model := ui.NewProgressModel()

	// Set up progress callback
	p := tea.NewProgram(model)
	c.SetProgressCallback(func(url string, status types.LinkStatus) {
		p.Send(ui.ProgressMsg{URL: url, Status: status})
	})

	// Run checker in goroutine
	go func() {
		result, err := c.CheckURLs(context.Background(), urls)
		if err != nil {
			p.Quit()
			return
		}
		p.Send(ui.DoneMsg{Result: result})
	}()

	// Run UI
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running UI: %w", err)
	}

	// Extract result from final model
	if m, ok := finalModel.(ui.ProgressModel); ok && m.Result() != nil {
		return m.Result(), nil
	}

	return nil, fmt.Errorf("no result available")
}

func runWithoutUI(c *checker.Checker, urls []string) (*types.CheckResult, error) {
	// Set up simple progress callback
	if cfg.Get().Verbose {
		c.SetProgressCallback(func(url string, status types.LinkStatus) {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", status, url)
		})
	}

	return c.CheckURLs(context.Background(), urls)
}

func outputResults(result *types.CheckResult) error {
	formatter := output.GetFormatter(cfg.Get().OutputFormat)

	// Determine output destination
	w := os.Stdout
	if cfg.Get().OutputFile != "" {
		f, err := os.Create(cfg.Get().OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return formatter.Format(result, w)
}

// executeCheck performs the actual link checking
func executeCheck(args []string) error {
	// Collect URLs to check
	urls, err := collectURLs(args)
	if err != nil {
		return fmt.Errorf("failed to collect URLs: %w", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs provided. Use --stdin to read from stdin or provide URLs as arguments")
	}

	// Create checker
	c, err := checker.New(cfg.Get())
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

	// Set up UI
	var result *types.CheckResult
	if cfg.Get().ShowProgress {
		result, err = runWithUI(c, urls)
	} else {
		result, err = runWithoutUI(c, urls)
	}

	if err != nil {
		return fmt.Errorf("check failed: %w", err)
	}

	// Format and output results
	if err := outputResults(result); err != nil {
		return fmt.Errorf("failed to output results: %w", err)
	}

	// Exit with error code if issues found
	if result.TotalDead > 0 || result.TotalErrors > 0 {
		os.Exit(1)
	}

	return nil
}

// runCheckCommand handles the check subcommand
func runCheckCommand(cmd *cobra.Command, args []string) error {
	// Force single mode
	cfg.Set("mode", types.ModeSingle)

	// Apply flags
	applyCheckFlags(cmd)

	return executeCheck(args)
}

// runCrawlCommand handles the crawl subcommand
func runCrawlCommand(cmd *cobra.Command, args []string) error {
	// Force crawler mode
	cfg.Set("mode", types.ModeCrawler)

	// Apply flags
	applyCrawlFlags(cmd)

	return executeCheck(args)
}

// applyCheckFlags applies check command flags to config
func applyCheckFlags(cmd *cobra.Command) {
	if cmd.Flags().Changed("concurrency") {
		cfg.Set("concurrency", flagConcurrency)
	}
	if cmd.Flags().Changed("timeout") {
		cfg.Set("timeout", flagTimeout)
	}
	if cmd.Flags().Changed("verbose") {
		cfg.Set("verbose", flagVerbose)
	}
	if cmd.Flags().Changed("no-progress") {
		cfg.Set("show_progress", !flagNoProgress)
	}
	if cmd.Flags().Changed("output-format") {
		cfg.Set("output_format", types.OutputFormat(flagOutputFormat))
	}
	if cmd.Flags().Changed("output-file") {
		cfg.Set("output_file", flagOutputFile)
	}
}

// applyCrawlFlags applies crawl command flags to config
func applyCrawlFlags(cmd *cobra.Command) {
	if cmd.Flags().Changed("max-depth") {
		cfg.Set("max_depth", flagMaxDepth)
	}
	if cmd.Flags().Changed("concurrency") {
		cfg.Set("concurrency", flagConcurrency)
	}
	if cmd.Flags().Changed("timeout") {
		cfg.Set("timeout", flagTimeout)
	}
	if cmd.Flags().Changed("verbose") {
		cfg.Set("verbose", flagVerbose)
	}
	if cmd.Flags().Changed("no-progress") {
		cfg.Set("show_progress", !flagNoProgress)
	}
	if cmd.Flags().Changed("output-format") {
		cfg.Set("output_format", types.OutputFormat(flagOutputFormat))
	}
	if cmd.Flags().Changed("output-file") {
		cfg.Set("output_file", flagOutputFile)
	}
}
