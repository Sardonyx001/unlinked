# Unlinked

A modern, fast, and beautiful dead link checker CLI built with Go. Unlinked helps you find and fix broken links on websites with style.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Features

- **Fast & Concurrent** - Check multiple links simultaneously with configurable concurrency
- **Beautiful UI** - Interactive terminal UI powered by Bubble Tea with real-time progress
- **Two Modes**:
  - **Single Mode** - Check specific URLs directly
  - **Crawler Mode** - Discover and check all links on a website
- **Multiple Output Formats** - Plaintext, Markdown, HTML, and JSON
- **Highly Configurable** - YAML configuration with CLI flags and environment variables
- **Smart Filtering** - Ignore patterns, domain restrictions, and robots.txt support
- **Detailed Reports** - Comprehensive statistics and link analysis
- **Redirect Handling** - Track and report HTTP redirects
- **Timeout Control** - Configurable timeouts and retry logic
- **Stdin Support** - Pipe URLs from other tools or files

## Installation

### Using Go Install

```bash
go install github.com/sardonyx001/unlinked/cmd/unlinked@latest
```

### From Source

```bash
git clone https://github.com/sardonyx001/unlinked.git
cd unlinked
task build
# Binary will be in bin/unlinked
```

### Using Task

Unlinked uses [Task](https://taskfile.dev/) for build automation. Install Task first:

```bash
# macOS
brew install go-task

# Go
go install github.com/go-task/task/v3/cmd/task@latest
```

Then build:

```bash
task build
```

## Quick Start

### Check a Single URL

```bash
unlinked https://example.com
```

### Check Multiple URLs

```bash
unlinked https://example.com https://golang.org https://github.com
```

### Crawler Mode - Find All Links

```bash
unlinked --mode=crawler https://example.com
```

### Read URLs from File

```bash
cat urls.txt | unlinked --stdin
```

### Output as Markdown

```bash
unlinked --output-format=markdown --output-file=report.md https://example.com
```

## Usage

```
unlinked [flags] [urls...]

Commands:
  version                        Print version information

Flags:
  -m, --mode string              Check mode: single or crawler (default "single")
  -c, --concurrency int          Number of concurrent checks (default 10)
  -t, --timeout int              Timeout in seconds for each request (default 30)
      --max-depth int            Maximum crawl depth (crawler mode) (default 3)
  -f, --output-format string     Output format: plaintext, markdown, html, json (default "plaintext")
  -o, --output-file string       Output file (default stdout)
  -v, --verbose                  Verbose output
      --no-progress              Disable progress display
      --stdin                    Read URLs from stdin
      --config string            Config file (default ~/.config/unlinked/config.yaml)
  -h, --help                     Help for unlinked

Examples:
  # Check version
  unlinked version
  unlinked version --short
```

## Configuration

Unlinked supports configuration via:

1. **Config file** (`~/.config/unlinked/config.yaml` or `./config.yaml`)
2. **Environment variables** (prefixed with `UNLINKED_`)
3. **Command-line flags** (highest priority)

### Example Configuration File

Create `~/.config/unlinked/config.yaml`:

```yaml
# Check mode: single or crawler
mode: single

# Output settings
output_format: plaintext
output_file: ""

# Performance settings
concurrency: 10
timeout: 30  # seconds

# Crawler settings
max_depth: 3
follow_redirects: true
check_external_only: false

# Behavior settings
respect_robots_txt: true
user_agent: "Unlinked/1.0 (Dead Link Checker)"

# Domain restrictions (crawler mode)
allowed_domains:
  - example.com
  - subdomain.example.com

# Patterns to ignore (regex)
ignore_patterns:
  - ".*\\.pdf$"
  - ".*\\.zip$"
  - "#.*"  # Anchors
  - "mailto:.*"  # Email links

# Display settings
verbose: false
show_progress: true
```

### Environment Variables

```bash
export UNLINKED_MODE=crawler
export UNLINKED_CONCURRENCY=20
export UNLINKED_TIMEOUT=60
export UNLINKED_OUTPUT_FORMAT=markdown
```

## Examples

### Basic Examples

```bash
# Check single URL with progress
unlinked https://example.com

# Check without progress (for scripts)
unlinked --no-progress https://example.com

# Check with verbose output
unlinked -v https://example.com
```

### Crawler Examples

```bash
# Crawl and check all links (depth 3)
unlinked --mode=crawler https://example.com

# Shallow crawl (depth 1)
unlinked --mode=crawler --max-depth=1 https://example.com

# Fast crawling with high concurrency
unlinked --mode=crawler --concurrency=50 https://example.com
```

### Output Format Examples

```bash
# Generate Markdown report
unlinked --output-format=markdown --output-file=report.md https://example.com

# Generate HTML report
unlinked --output-format=html --output-file=report.html https://example.com

# JSON output for programmatic use
unlinked --output-format=json --output-file=report.json https://example.com

# Pretty print to terminal
unlinked --output-format=plaintext https://example.com
```

### Pipeline Examples

```bash
# Check URLs from a file
cat urls.txt | unlinked --stdin

# Check URLs from another command
grep -r "http" docs/ | cut -d'"' -f2 | unlinked --stdin

# Check and save results
cat urls.txt | unlinked --stdin --output-format=markdown -o report.md

# Check multiple sites from CSV
cut -d',' -f1 sites.csv | tail -n +2 | unlinked --stdin
```

### Advanced Examples

```bash
# Custom timeout and concurrency
unlinked --timeout=60 --concurrency=20 https://example.com

# With custom config
unlinked --config=my-config.yaml https://example.com

# Crawler with domain restrictions (via config)
unlinked --mode=crawler --config=crawler-config.yaml https://example.com

# Quiet mode for CI/CD
unlinked --no-progress --output-format=json https://example.com
```

## Output Formats

### Plaintext

Simple, readable output perfect for terminal viewing:

```
Link Check Report
=================

Summary:
  Start Time:    2024-01-15T10:30:00Z
  End Time:      2024-01-15T10:30:05Z
  Duration:      5.2s
  Total Checked: 45
  OK:            42
  Dead:          2
  Redirects:     1
  Errors:        0

Dead Links (2):
----------------
  [404] https://example.com/missing
       Found on: https://example.com/index.html
...
```

### Markdown

GitHub-compatible Markdown with emojis and tables:

```markdown
# Link Check Report

## Summary

| Metric | Value |
|--------|-------|
| Total Checked | 45 |
| ✅ OK | 42 |
| ❌ Dead | 2 |
| 🔀 Redirects | 1 |

## ❌ Dead Links (2)

- **[404]** `https://example.com/missing`
  - Found on: <https://example.com/index.html>
...
```

### HTML

Beautiful, styled HTML report:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Link Check Report</title>
    <!-- Includes responsive CSS styling -->
</head>
<body>
    <!-- Interactive report with color-coded results -->
</body>
</html>
```

### JSON

Structured data for programmatic processing:

```json
{
  "start_time": "2024-01-15T10:30:00Z",
  "end_time": "2024-01-15T10:30:05Z",
  "total_checked": 45,
  "total_ok": 42,
  "total_dead": 2,
  "links": [...]
}
```

## Development

### Prerequisites

- Go 1.25 or higher
- Task (optional, for build automation)

### Building

```bash
# Build for current platform
task build

# Build for all platforms
task build-all

# Install locally
task install
```

### Testing

```bash
# Run tests
task test

# Run tests with coverage
task test-coverage

# Run benchmarks
task bench
```

### Linting

```bash
# Format and vet
task lint

# Tidy dependencies
task tidy
```

### Available Tasks

```bash
# List all available tasks
task --list

# or
task
```

Common tasks:

- `task build` - Build the binary
- `task test` - Run tests
- `task lint` - Run linters
- `task clean` - Clean build artifacts
- `task install` - Install to $GOPATH/bin
- `task run -- [args]` - Run with arguments
- `task example-crawler` - Run crawler example

## Project Structure

```
unlinked/
├── cmd/
│   └── unlinked/          # CLI entry point
│       ├── main.go
│       └── root.go        # Cobra root command
├── internal/
│   ├── checker/           # Link checking engine
│   │   └── checker.go
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── output/            # Output formatters
│   │   └── formatter.go
│   └── ui/                # Terminal UI (Bubble Tea)
│       └── progress.go
├── pkg/
│   └── types/             # Shared types
│       └── types.go
├── Taskfile.yml           # Build automation
├── config.example.yaml    # Example configuration
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Built with amazing Go libraries:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output
- [Colly](https://github.com/gocolly/colly) - Web scraping framework

## Contact

Issues and questions can be posted on the [GitHub Issues](https://github.com/sardonyx001/unlinked/issues) page.

---

Made with Go
