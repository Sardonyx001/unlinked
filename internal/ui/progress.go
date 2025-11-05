package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sardonyx001/unlinked/pkg/types"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	statusOKStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	statusDeadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6C4D")).
			Bold(true)

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F59E0B"))

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)

	statsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginTop(1)
)

// ProgressMsg represents a progress update
type ProgressMsg struct {
	URL    string
	Status types.LinkStatus
}

// DoneMsg signals completion
type DoneMsg struct {
	Result *types.CheckResult
}

// ProgressModel is the Bubble Tea model for progress display
type ProgressModel struct {
	spinner      spinner.Model
	progress     progress.Model
	currentURL   string
	currentStatus types.LinkStatus
	stats        Stats
	width        int
	done         bool
	result       *types.CheckResult
}

// Stats holds checking statistics
type Stats struct {
	Total     int
	OK        int
	Dead      int
	Redirects int
	Errors    int
}

// NewProgressModel creates a new progress model
func NewProgressModel() ProgressModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	p := progress.New(progress.WithDefaultGradient())
	p.Width = 40

	return ProgressModel{
		spinner:  s,
		progress: p,
		stats:    Stats{},
		width:    80,
	}
}

// Init initializes the model
func (m ProgressModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case ProgressMsg:
		m.currentURL = msg.URL
		m.currentStatus = msg.Status
		m.stats.Total++
		switch msg.Status {
		case types.StatusOK:
			m.stats.OK++
		case types.StatusDead:
			m.stats.Dead++
		case types.StatusRedirect:
			m.stats.Redirects++
		case types.StatusError, types.StatusTimeout:
			m.stats.Errors++
		}
		return m, nil

	case DoneMsg:
		m.done = true
		m.result = msg.Result
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m ProgressModel) View() string {
	if m.done && m.result != nil {
		return m.renderFinalReport()
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🔗 Unlinked - Dead Link Checker"))
	b.WriteString("\n\n")

	// Current status
	if m.currentURL != "" {
		b.WriteString(m.spinner.View())
		b.WriteString(" Checking: ")
		b.WriteString(urlStyle.Render(truncateURL(m.currentURL, m.width-20)))
		b.WriteString(" ")
		b.WriteString(m.renderStatus(m.currentStatus))
		b.WriteString("\n\n")
	}

	// Statistics
	b.WriteString(m.renderStats())
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("Press q or Ctrl+C to quit"))

	return b.String()
}

// renderStatus renders a status badge
func (m ProgressModel) renderStatus(status types.LinkStatus) string {
	switch status {
	case types.StatusOK:
		return statusOKStyle.Render("✓ OK")
	case types.StatusDead:
		return statusDeadStyle.Render("✗ DEAD")
	case types.StatusRedirect:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")).Render("↻ REDIRECT")
	case types.StatusError, types.StatusTimeout:
		return statusErrorStyle.Render("! ERROR")
	default:
		return ""
	}
}

// renderStats renders the statistics box
func (m ProgressModel) renderStats() string {
	return statsStyle.Render(fmt.Sprintf(
		"Total: %d  |  %s: %d  |  %s: %d  |  %s: %d  |  %s: %d",
		m.stats.Total,
		statusOKStyle.Render("OK"),
		m.stats.OK,
		statusDeadStyle.Render("Dead"),
		m.stats.Dead,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")).Render("Redirects"),
		m.stats.Redirects,
		statusErrorStyle.Render("Errors"),
		m.stats.Errors,
	))
}

// renderFinalReport renders the final report summary
func (m ProgressModel) renderFinalReport() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🎉 Check Complete!"))
	b.WriteString("\n\n")

	result := m.result
	b.WriteString(fmt.Sprintf("Duration: %s\n", result.Duration.Round(1)))
	b.WriteString(fmt.Sprintf("Total Checked: %d\n\n", result.TotalChecked))

	b.WriteString(statsStyle.Render(fmt.Sprintf(
		"%s: %d  |  %s: %d  |  %s: %d  |  %s: %d",
		statusOKStyle.Render("✓ OK"),
		result.TotalOK,
		statusDeadStyle.Render("✗ Dead"),
		result.TotalDead,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")).Render("↻ Redirects"),
		result.TotalRedirect,
		statusErrorStyle.Render("! Errors"),
		result.TotalErrors,
	)))

	if result.TotalDead > 0 || result.TotalErrors > 0 {
		b.WriteString("\n\n")
		b.WriteString(statusErrorStyle.Render(fmt.Sprintf(
			"⚠️  Found %d issues that need attention",
			result.TotalDead+result.TotalErrors,
		)))
	} else {
		b.WriteString("\n\n")
		b.WriteString(statusOKStyle.Render("✨ All links are healthy!"))
	}

	return b.String()
}

// truncateURL truncates a URL to fit within a given width
func truncateURL(url string, maxWidth int) string {
	if len(url) <= maxWidth {
		return url
	}
	return url[:maxWidth-3] + "..."
}

// Result returns the final result
func (m ProgressModel) Result() *types.CheckResult {
	return m.result
}
