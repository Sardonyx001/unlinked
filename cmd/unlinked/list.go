package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	listVerbose bool
	listStyles  = struct {
		Title       lipgloss.Style
		Command     lipgloss.Style
		Description lipgloss.Style
		Category    lipgloss.Style
	}{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1),
		Command: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")),
		Category: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F59E0B")).
			MarginTop(1).
			MarginBottom(0),
	}
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available commands",
	Long:  `Display a formatted list of all available commands with descriptions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if listVerbose {
			showVerboseList()
		} else {
			showCompactList()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "show detailed command information")
}

func showCompactList() {
	fmt.Println(listStyles.Title.Render("Unlinked - Available Commands"))
	fmt.Println()

	commands := []struct {
		name string
		desc string
	}{
		{"check", "Check specific URLs for broken links"},
		{"crawl", "Crawl a website and check all discovered links"},
		{"version", "Print version information"},
		{"completion", "Generate shell completion scripts"},
		{"help", "Show help for any command"},
		{"list", "List all available commands"},
	}

	maxLen := 0
	for _, cmd := range commands {
		if len(cmd.name) > maxLen {
			maxLen = len(cmd.name)
		}
	}

	for _, cmd := range commands {
		padding := strings.Repeat(" ", maxLen-len(cmd.name)+2)
		fmt.Printf("  %s%s%s\n",
			listStyles.Command.Render(cmd.name),
			padding,
			listStyles.Description.Render(cmd.desc),
		)
	}
	fmt.Println()
	fmt.Println("Run 'unlinked <command> --help' for more information on a command.")
}

func showVerboseList() {
	fmt.Println(listStyles.Title.Render("Unlinked - Command Reference"))
	fmt.Println()

	// Core Commands
	fmt.Println(listStyles.Category.Render("Core Commands:"))
	fmt.Printf("  %s\n", listStyles.Command.Render("check [urls...]"))
	fmt.Printf("    %s\n", listStyles.Description.Render("Check one or more specific URLs for accessibility"))
	fmt.Printf("    %s\n\n", listStyles.Description.Render("Example: unlinked check https://example.com"))

	fmt.Printf("  %s\n", listStyles.Command.Render("crawl [url]"))
	fmt.Printf("    %s\n", listStyles.Description.Render("Crawl a website and check all discovered links"))
	fmt.Printf("    %s\n\n", listStyles.Description.Render("Example: unlinked crawl --max-depth=2 https://example.com"))

	// Utility Commands
	fmt.Println(listStyles.Category.Render("Utility Commands:"))
	fmt.Printf("  %s\n", listStyles.Command.Render("version [--short]"))
	fmt.Printf("    %s\n", listStyles.Description.Render("Display version and build information"))
	fmt.Printf("    %s\n\n", listStyles.Description.Render("Example: unlinked version --short"))

	fmt.Printf("  %s\n", listStyles.Command.Render("completion [bash|zsh|fish|powershell]"))
	fmt.Printf("    %s\n", listStyles.Description.Render("Generate shell completion scripts"))
	fmt.Printf("    %s\n\n", listStyles.Description.Render("Example: source <(unlinked completion bash)"))

	fmt.Printf("  %s\n", listStyles.Command.Render("help [command]"))
	fmt.Printf("    %s\n", listStyles.Description.Render("Show help information for any command"))
	fmt.Printf("    %s\n\n", listStyles.Description.Render("Example: unlinked help check"))

	fmt.Printf("  %s\n", listStyles.Command.Render("list [--verbose]"))
	fmt.Printf("    %s\n", listStyles.Description.Render("List all available commands"))
	fmt.Printf("    %s\n\n", listStyles.Description.Render("Example: unlinked list -v"))

	fmt.Println(listStyles.Category.Render("Common Flags:"))
	fmt.Printf("  %s\n", listStyles.Description.Render("--config          Config file path"))
	fmt.Printf("  %s\n", listStyles.Description.Render("-f, --output-format   Output format (plaintext, markdown, html, json)"))
	fmt.Printf("  %s\n", listStyles.Description.Render("-o, --output-file     Output file path"))
	fmt.Printf("  %s\n", listStyles.Description.Render("-c, --concurrency     Number of concurrent checks"))
	fmt.Printf("  %s\n", listStyles.Description.Render("-t, --timeout         Request timeout in seconds"))
	fmt.Printf("  %s\n", listStyles.Description.Render("-v, --verbose         Verbose output"))
	fmt.Printf("  %s\n", listStyles.Description.Render("--no-progress        Disable progress display"))
	fmt.Println()
}
