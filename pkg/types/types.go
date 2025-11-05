package types

import "time"

// CheckMode defines how URLs should be checked
type CheckMode string

const (
	// ModeSingle checks only the provided URL(s)
	ModeSingle CheckMode = "single"
	// ModeCrawler crawls the URL and checks all discovered links
	ModeCrawler CheckMode = "crawler"
)

// OutputFormat defines the output format for results
type OutputFormat string

const (
	FormatPlaintext OutputFormat = "plaintext"
	FormatMarkdown  OutputFormat = "markdown"
	FormatHTML      OutputFormat = "html"
	FormatJSON      OutputFormat = "json"
)

// LinkStatus represents the status of a checked link
type LinkStatus string

const (
	StatusOK       LinkStatus = "ok"
	StatusDead     LinkStatus = "dead"
	StatusRedirect LinkStatus = "redirect"
	StatusTimeout  LinkStatus = "timeout"
	StatusError    LinkStatus = "error"
	StatusSkipped  LinkStatus = "skipped"
)

// LinkResult represents the result of checking a single link
type LinkResult struct {
	URL           string        `json:"url"`
	Status        LinkStatus    `json:"status"`
	StatusCode    int           `json:"status_code"`
	Error         string        `json:"error,omitempty"`
	RedirectURL   string        `json:"redirect_url,omitempty"`
	FoundOn       string        `json:"found_on,omitempty"` // Parent URL where link was found
	ResponseTime  time.Duration `json:"response_time"`
	CheckedAt     time.Time     `json:"checked_at"`
	ContentType   string        `json:"content_type,omitempty"`
	ContentLength int64         `json:"content_length,omitempty"`
}

// CheckResult represents the complete result of a check operation
type CheckResult struct {
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	TotalChecked  int           `json:"total_checked"`
	TotalOK       int           `json:"total_ok"`
	TotalDead     int           `json:"total_dead"`
	TotalRedirect int           `json:"total_redirect"`
	TotalErrors   int           `json:"total_errors"`
	Links         []LinkResult  `json:"links"`
	Duration      time.Duration `json:"duration"`
}

// Config represents the application configuration
type Config struct {
	Mode              CheckMode    `mapstructure:"mode"`
	OutputFormat      OutputFormat `mapstructure:"output_format"`
	OutputFile        string       `mapstructure:"output_file"`
	Concurrency       int          `mapstructure:"concurrency"`
	Timeout           int          `mapstructure:"timeout"` // in seconds
	MaxDepth          int          `mapstructure:"max_depth"`
	FollowRedirects   bool         `mapstructure:"follow_redirects"`
	CheckExternalOnly bool         `mapstructure:"check_external_only"`
	UserAgent         string       `mapstructure:"user_agent"`
	RespectRobotsTxt  bool         `mapstructure:"respect_robots_txt"`
	AllowedDomains    []string     `mapstructure:"allowed_domains"`
	IgnorePatterns    []string     `mapstructure:"ignore_patterns"`
	Verbose           bool         `mapstructure:"verbose"`
	ShowProgress      bool         `mapstructure:"show_progress"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Mode:             ModeSingle,
		OutputFormat:     FormatPlaintext,
		Concurrency:      10,
		Timeout:          30,
		MaxDepth:         3,
		FollowRedirects:  true,
		UserAgent:        "Unlinked/1.0 (Dead Link Checker)",
		RespectRobotsTxt: true,
		Verbose:          false,
		ShowProgress:     true,
	}
}
