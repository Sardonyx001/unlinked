package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sardonyx001/unlinked/pkg/types"
)

// Formatter handles output formatting
type Formatter interface {
	Format(result *types.CheckResult, w io.Writer) error
}

// GetFormatter returns the appropriate formatter based on the output format
func GetFormatter(format types.OutputFormat) Formatter {
	switch format {
	case types.FormatMarkdown:
		return &MarkdownFormatter{}
	case types.FormatHTML:
		return &HTMLFormatter{}
	case types.FormatJSON:
		return &JSONFormatter{}
	default:
		return &PlaintextFormatter{}
	}
}

// PlaintextFormatter formats output as plain text
type PlaintextFormatter struct{}

func (f *PlaintextFormatter) Format(result *types.CheckResult, w io.Writer) error {
	fmt.Fprintf(w, "Link Check Report\n")
	fmt.Fprintf(w, "=================\n\n")
	fmt.Fprintf(w, "Summary:\n")
	fmt.Fprintf(w, "  Start Time:    %s\n", result.StartTime.Format(time.RFC3339))
	fmt.Fprintf(w, "  End Time:      %s\n", result.EndTime.Format(time.RFC3339))
	fmt.Fprintf(w, "  Duration:      %s\n", result.Duration.Round(time.Millisecond))
	fmt.Fprintf(w, "  Total Checked: %d\n", result.TotalChecked)
	fmt.Fprintf(w, "  OK:            %d\n", result.TotalOK)
	fmt.Fprintf(w, "  Dead:          %d\n", result.TotalDead)
	fmt.Fprintf(w, "  Redirects:     %d\n", result.TotalRedirect)
	fmt.Fprintf(w, "  Errors:        %d\n\n", result.TotalErrors)

	// Group links by status
	byStatus := groupByStatus(result.Links)

	if len(byStatus[types.StatusDead]) > 0 {
		fmt.Fprintf(w, "Dead Links (%d):\n", len(byStatus[types.StatusDead]))
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 80))
		for _, link := range byStatus[types.StatusDead] {
			fmt.Fprintf(w, "  [%d] %s\n", link.StatusCode, link.URL)
			if link.FoundOn != "" {
				fmt.Fprintf(w, "       Found on: %s\n", link.FoundOn)
			}
			if link.Error != "" {
				fmt.Fprintf(w, "       Error: %s\n", link.Error)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	if len(byStatus[types.StatusError]) > 0 || len(byStatus[types.StatusTimeout]) > 0 {
		errors := append(byStatus[types.StatusError], byStatus[types.StatusTimeout]...)
		fmt.Fprintf(w, "Errors (%d):\n", len(errors))
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 80))
		for _, link := range errors {
			fmt.Fprintf(w, "  [%s] %s\n", link.Status, link.URL)
			if link.FoundOn != "" {
				fmt.Fprintf(w, "       Found on: %s\n", link.FoundOn)
			}
			if link.Error != "" {
				fmt.Fprintf(w, "       Error: %s\n", link.Error)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	if len(byStatus[types.StatusRedirect]) > 0 {
		fmt.Fprintf(w, "Redirects (%d):\n", len(byStatus[types.StatusRedirect]))
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 80))
		for _, link := range byStatus[types.StatusRedirect] {
			fmt.Fprintf(w, "  [%d] %s\n", link.StatusCode, link.URL)
			if link.RedirectURL != "" {
				fmt.Fprintf(w, "       -> %s\n", link.RedirectURL)
			}
			if link.FoundOn != "" {
				fmt.Fprintf(w, "       Found on: %s\n", link.FoundOn)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}

// MarkdownFormatter formats output as Markdown
type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Format(result *types.CheckResult, w io.Writer) error {
	fmt.Fprintf(w, "# Link Check Report\n\n")

	// Summary table
	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "| Metric | Value |\n")
	fmt.Fprintf(w, "|--------|-------|\n")
	fmt.Fprintf(w, "| Start Time | %s |\n", result.StartTime.Format(time.RFC3339))
	fmt.Fprintf(w, "| End Time | %s |\n", result.EndTime.Format(time.RFC3339))
	fmt.Fprintf(w, "| Duration | %s |\n", result.Duration.Round(time.Millisecond))
	fmt.Fprintf(w, "| Total Checked | %d |\n", result.TotalChecked)
	fmt.Fprintf(w, "| ✅ OK | %d |\n", result.TotalOK)
	fmt.Fprintf(w, "| ❌ Dead | %d |\n", result.TotalDead)
	fmt.Fprintf(w, "| 🔀 Redirects | %d |\n", result.TotalRedirect)
	fmt.Fprintf(w, "| ⚠️ Errors | %d |\n\n", result.TotalErrors)

	// Group links by status
	byStatus := groupByStatus(result.Links)

	if len(byStatus[types.StatusDead]) > 0 {
		fmt.Fprintf(w, "## ❌ Dead Links (%d)\n\n", len(byStatus[types.StatusDead]))
		for _, link := range byStatus[types.StatusDead] {
			fmt.Fprintf(w, "- **[%d]** `%s`\n", link.StatusCode, link.URL)
			if link.FoundOn != "" {
				fmt.Fprintf(w, "  - Found on: <%s>\n", link.FoundOn)
			}
			if link.Error != "" {
				fmt.Fprintf(w, "  - Error: `%s`\n", link.Error)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	if len(byStatus[types.StatusError]) > 0 || len(byStatus[types.StatusTimeout]) > 0 {
		errors := append(byStatus[types.StatusError], byStatus[types.StatusTimeout]...)
		fmt.Fprintf(w, "## ⚠️ Errors (%d)\n\n", len(errors))
		for _, link := range errors {
			fmt.Fprintf(w, "- **[%s]** `%s`\n", link.Status, link.URL)
			if link.FoundOn != "" {
				fmt.Fprintf(w, "  - Found on: <%s>\n", link.FoundOn)
			}
			if link.Error != "" {
				fmt.Fprintf(w, "  - Error: `%s`\n", link.Error)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	if len(byStatus[types.StatusRedirect]) > 0 {
		fmt.Fprintf(w, "## 🔀 Redirects (%d)\n\n", len(byStatus[types.StatusRedirect]))
		for _, link := range byStatus[types.StatusRedirect] {
			fmt.Fprintf(w, "- **[%d]** `%s`\n", link.StatusCode, link.URL)
			if link.RedirectURL != "" {
				fmt.Fprintf(w, "  - Redirects to: <%s>\n", link.RedirectURL)
			}
			if link.FoundOn != "" {
				fmt.Fprintf(w, "  - Found on: <%s>\n", link.FoundOn)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}

// HTMLFormatter formats output as HTML
type HTMLFormatter struct{}

func (f *HTMLFormatter) Format(result *types.CheckResult, w io.Writer) error {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Link Check Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 8px;
            padding: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            border-bottom: 3px solid #4CAF50;
            padding-bottom: 10px;
        }
        h2 {
            color: #555;
            margin-top: 30px;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }
        .stat-card {
            background: #f9f9f9;
            padding: 15px;
            border-radius: 6px;
            border-left: 4px solid #4CAF50;
        }
        .stat-card.dead { border-left-color: #f44336; }
        .stat-card.error { border-left-color: #ff9800; }
        .stat-card.redirect { border-left-color: #2196F3; }
        .stat-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
        }
        .stat-value {
            font-size: 28px;
            font-weight: bold;
            color: #333;
        }
        .link-item {
            background: #f9f9f9;
            margin: 10px 0;
            padding: 15px;
            border-radius: 6px;
            border-left: 4px solid #ccc;
        }
        .link-item.dead { border-left-color: #f44336; }
        .link-item.error { border-left-color: #ff9800; }
        .link-item.redirect { border-left-color: #2196F3; }
        .link-url {
            font-family: monospace;
            word-break: break-all;
            font-weight: bold;
        }
        .link-meta {
            font-size: 13px;
            color: #666;
            margin-top: 8px;
        }
        .badge {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 12px;
            font-weight: bold;
            margin-right: 5px;
        }
        .badge.dead { background: #f44336; color: white; }
        .badge.error { background: #ff9800; color: white; }
        .badge.redirect { background: #2196F3; color: white; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔗 Link Check Report</h1>

        <div class="summary">
            <div class="stat-card">
                <div class="stat-label">Total Checked</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">✅ OK</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card dead">
                <div class="stat-label">❌ Dead</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card redirect">
                <div class="stat-label">🔀 Redirects</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card error">
                <div class="stat-label">⚠️ Errors</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Duration</div>
                <div class="stat-value" style="font-size: 18px;">%s</div>
            </div>
        </div>
`, result.TotalChecked, result.TotalOK, result.TotalDead, result.TotalRedirect,
   result.TotalErrors, result.Duration.Round(time.Millisecond))

	byStatus := groupByStatus(result.Links)

	if len(byStatus[types.StatusDead]) > 0 {
		fmt.Fprintf(w, "<h2>❌ Dead Links (%d)</h2>\n", len(byStatus[types.StatusDead]))
		for _, link := range byStatus[types.StatusDead] {
			fmt.Fprintf(w, `<div class="link-item dead">
    <div><span class="badge dead">%d</span><span class="link-url">%s</span></div>
`, link.StatusCode, escapeHTML(link.URL))
			if link.FoundOn != "" {
				fmt.Fprintf(w, `    <div class="link-meta">Found on: <a href="%s">%s</a></div>
`, link.FoundOn, escapeHTML(link.FoundOn))
			}
			if link.Error != "" {
				fmt.Fprintf(w, `    <div class="link-meta">Error: %s</div>
`, escapeHTML(link.Error))
			}
			fmt.Fprintf(w, "</div>\n")
		}
	}

	if len(byStatus[types.StatusError]) > 0 || len(byStatus[types.StatusTimeout]) > 0 {
		errors := append(byStatus[types.StatusError], byStatus[types.StatusTimeout]...)
		fmt.Fprintf(w, "<h2>⚠️ Errors (%d)</h2>\n", len(errors))
		for _, link := range errors {
			fmt.Fprintf(w, `<div class="link-item error">
    <div><span class="badge error">%s</span><span class="link-url">%s</span></div>
`, link.Status, escapeHTML(link.URL))
			if link.FoundOn != "" {
				fmt.Fprintf(w, `    <div class="link-meta">Found on: <a href="%s">%s</a></div>
`, link.FoundOn, escapeHTML(link.FoundOn))
			}
			if link.Error != "" {
				fmt.Fprintf(w, `    <div class="link-meta">Error: %s</div>
`, escapeHTML(link.Error))
			}
			fmt.Fprintf(w, "</div>\n")
		}
	}

	if len(byStatus[types.StatusRedirect]) > 0 {
		fmt.Fprintf(w, "<h2>🔀 Redirects (%d)</h2>\n", len(byStatus[types.StatusRedirect]))
		for _, link := range byStatus[types.StatusRedirect] {
			fmt.Fprintf(w, `<div class="link-item redirect">
    <div><span class="badge redirect">%d</span><span class="link-url">%s</span></div>
`, link.StatusCode, escapeHTML(link.URL))
			if link.RedirectURL != "" {
				fmt.Fprintf(w, `    <div class="link-meta">Redirects to: <a href="%s">%s</a></div>
`, link.RedirectURL, escapeHTML(link.RedirectURL))
			}
			if link.FoundOn != "" {
				fmt.Fprintf(w, `    <div class="link-meta">Found on: <a href="%s">%s</a></div>
`, link.FoundOn, escapeHTML(link.FoundOn))
			}
			fmt.Fprintf(w, "</div>\n")
		}
	}

	fmt.Fprintf(w, `
    </div>
</body>
</html>
`)
	return nil
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

func (f *JSONFormatter) Format(result *types.CheckResult, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// Helper functions

func groupByStatus(links []types.LinkResult) map[types.LinkStatus][]types.LinkResult {
	grouped := make(map[types.LinkStatus][]types.LinkResult)
	for _, link := range links {
		grouped[link.Status] = append(grouped[link.Status], link)
	}
	return grouped
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
