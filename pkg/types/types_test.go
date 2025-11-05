package types

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Test default values
	if config.Mode != ModeSingle {
		t.Errorf("Expected default mode to be %s, got %s", ModeSingle, config.Mode)
	}

	if config.OutputFormat != FormatPlaintext {
		t.Errorf("Expected default output format to be %s, got %s", FormatPlaintext, config.OutputFormat)
	}

	if config.Concurrency != 10 {
		t.Errorf("Expected default concurrency to be 10, got %d", config.Concurrency)
	}

	if config.Timeout != 30 {
		t.Errorf("Expected default timeout to be 30, got %d", config.Timeout)
	}

	if config.MaxDepth != 3 {
		t.Errorf("Expected default max depth to be 3, got %d", config.MaxDepth)
	}

	if !config.FollowRedirects {
		t.Error("Expected default FollowRedirects to be true")
	}

	if !config.RespectRobotsTxt {
		t.Error("Expected default RespectRobotsTxt to be true")
	}

	if !config.ShowProgress {
		t.Error("Expected default ShowProgress to be true")
	}

	if config.Verbose {
		t.Error("Expected default Verbose to be false")
	}

	if config.UserAgent != "Unlinked/1.0 (Dead Link Checker)" {
		t.Errorf("Unexpected default user agent: %s", config.UserAgent)
	}
}

func TestLinkStatus(t *testing.T) {
	tests := []struct {
		name   string
		status LinkStatus
		value  string
	}{
		{"OK status", StatusOK, "ok"},
		{"Dead status", StatusDead, "dead"},
		{"Redirect status", StatusRedirect, "redirect"},
		{"Timeout status", StatusTimeout, "timeout"},
		{"Error status", StatusError, "error"},
		{"Skipped status", StatusSkipped, "skipped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.value {
				t.Errorf("Expected status value %s, got %s", tt.value, string(tt.status))
			}
		})
	}
}

func TestCheckMode(t *testing.T) {
	if ModeSingle != "single" {
		t.Errorf("Expected ModeSingle to be 'single', got '%s'", ModeSingle)
	}
	if ModeCrawler != "crawler" {
		t.Errorf("Expected ModeCrawler to be 'crawler', got '%s'", ModeCrawler)
	}
}

func TestOutputFormat(t *testing.T) {
	formats := map[OutputFormat]string{
		FormatPlaintext: "plaintext",
		FormatMarkdown:  "markdown",
		FormatHTML:      "html",
		FormatJSON:      "json",
	}

	for format, expected := range formats {
		if string(format) != expected {
			t.Errorf("Expected format value %s, got %s", expected, string(format))
		}
	}
}
