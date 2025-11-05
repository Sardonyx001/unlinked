package checker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/sardonyx001/unlinked/pkg/types"
)

// Checker handles link checking operations
type Checker struct {
	config      *types.Config
	results     []types.LinkResult
	visited     map[string]bool
	mu          sync.Mutex
	client      *http.Client
	onProgress  func(url string, status types.LinkStatus)
	ignoreRegex []*regexp.Regexp
}

// New creates a new link checker
func New(config *types.Config) (*Checker, error) {
	c := &Checker{
		config:  config,
		results: make([]types.LinkResult, 0),
		visited: make(map[string]bool),
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if !config.FollowRedirects {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		ignoreRegex: make([]*regexp.Regexp, 0),
	}

	// Compile ignore patterns
	for _, pattern := range config.IgnorePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore pattern %q: %w", pattern, err)
		}
		c.ignoreRegex = append(c.ignoreRegex, re)
	}

	return c, nil
}

// SetProgressCallback sets a callback for progress updates
func (c *Checker) SetProgressCallback(fn func(url string, status types.LinkStatus)) {
	c.onProgress = fn
}

// CheckURLs checks a list of URLs based on the configured mode
func (c *Checker) CheckURLs(ctx context.Context, urls []string) (*types.CheckResult, error) {
	startTime := time.Now()

	for _, u := range urls {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if c.config.Mode == types.ModeCrawler {
			if err := c.crawlAndCheck(ctx, u); err != nil {
				return nil, err
			}
		} else {
			c.checkSingleURL(u, "")
		}
	}

	endTime := time.Now()
	return c.buildResult(startTime, endTime), nil
}

// checkSingleURL checks a single URL
func (c *Checker) checkSingleURL(targetURL, foundOn string) types.LinkResult {
	c.mu.Lock()
	if c.visited[targetURL] {
		c.mu.Unlock()
		return types.LinkResult{URL: targetURL, Status: types.StatusSkipped}
	}
	c.visited[targetURL] = true
	c.mu.Unlock()

	// Check if URL should be ignored
	if c.shouldIgnore(targetURL) {
		result := types.LinkResult{
			URL:       targetURL,
			Status:    types.StatusSkipped,
			FoundOn:   foundOn,
			CheckedAt: time.Now(),
		}
		c.addResult(result)
		return result
	}

	startTime := time.Now()

	req, err := http.NewRequest("HEAD", targetURL, nil)
	if err != nil {
		result := types.LinkResult{
			URL:       targetURL,
			Status:    types.StatusError,
			Error:     err.Error(),
			FoundOn:   foundOn,
			CheckedAt: time.Now(),
		}
		c.addResult(result)
		c.notifyProgress(targetURL, types.StatusError)
		return result
	}

	req.Header.Set("User-Agent", c.config.UserAgent)

	resp, err := c.client.Do(req)
	responseTime := time.Since(startTime)

	if err != nil {
		status := types.StatusError
		if err, ok := err.(net.Error); ok && err.Timeout() {
			status = types.StatusTimeout
		}
		result := types.LinkResult{
			URL:          targetURL,
			Status:       status,
			Error:        err.Error(),
			FoundOn:      foundOn,
			ResponseTime: responseTime,
			CheckedAt:    time.Now(),
		}
		c.addResult(result)
		c.notifyProgress(targetURL, status)
		return result
	}
	defer resp.Body.Close()

	// Determine status
	status := c.determineStatus(resp.StatusCode)

	result := types.LinkResult{
		URL:           targetURL,
		Status:        status,
		StatusCode:    resp.StatusCode,
		FoundOn:       foundOn,
		ResponseTime:  responseTime,
		CheckedAt:     time.Now(),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}

	// Handle redirects
	if status == types.StatusRedirect && resp.Header.Get("Location") != "" {
		result.RedirectURL = resp.Header.Get("Location")
	}

	c.addResult(result)
	c.notifyProgress(targetURL, status)

	return result
}

// crawlAndCheck crawls a URL and checks all discovered links
func (c *Checker) crawlAndCheck(ctx context.Context, startURL string) error {
	collector := colly.NewCollector(
		colly.MaxDepth(c.config.MaxDepth),
		colly.Async(true),
		colly.UserAgent(c.config.UserAgent),
	)

	// Set allowed domains if specified
	if len(c.config.AllowedDomains) > 0 {
		collector.AllowedDomains = c.config.AllowedDomains
	}

	// Configure parallelism
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: c.config.Concurrency,
	})

	// Extract and check all links
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		// Check the link
		c.checkSingleURL(link, e.Request.URL.String())

		// Visit the link if in crawler mode (to find more links)
		if c.config.Mode == types.ModeCrawler {
			e.Request.Visit(link)
		}
	})

	// Handle errors
	collector.OnError(func(r *colly.Response, err error) {
		result := types.LinkResult{
			URL:        r.Request.URL.String(),
			Status:     types.StatusError,
			StatusCode: r.StatusCode,
			Error:      err.Error(),
			CheckedAt:  time.Now(),
		}
		c.addResult(result)
		c.notifyProgress(r.Request.URL.String(), types.StatusError)
	})

	// Start crawling
	if err := collector.Visit(startURL); err != nil {
		return fmt.Errorf("failed to start crawling: %w", err)
	}

	// Wait for all async requests to complete
	collector.Wait()

	return nil
}

// shouldIgnore checks if a URL should be ignored based on patterns
func (c *Checker) shouldIgnore(targetURL string) bool {
	for _, re := range c.ignoreRegex {
		if re.MatchString(targetURL) {
			return true
		}
	}
	return false
}

// determineStatus determines the link status based on HTTP status code
func (c *Checker) determineStatus(code int) types.LinkStatus {
	switch {
	case code >= 200 && code < 300:
		return types.StatusOK
	case code >= 300 && code < 400:
		return types.StatusRedirect
	default:
		return types.StatusDead
	}
}

// addResult adds a result to the results list (thread-safe)
func (c *Checker) addResult(result types.LinkResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = append(c.results, result)
}

// notifyProgress notifies the progress callback if set
func (c *Checker) notifyProgress(url string, status types.LinkStatus) {
	if c.onProgress != nil {
		c.onProgress(url, status)
	}
}

// buildResult builds the final check result
func (c *Checker) buildResult(startTime, endTime time.Time) *types.CheckResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := &types.CheckResult{
		StartTime:    startTime,
		EndTime:      endTime,
		Duration:     endTime.Sub(startTime),
		TotalChecked: len(c.results),
		Links:        c.results,
	}

	// Calculate statistics
	for _, link := range c.results {
		switch link.Status {
		case types.StatusOK:
			result.TotalOK++
		case types.StatusDead:
			result.TotalDead++
		case types.StatusRedirect:
			result.TotalRedirect++
		case types.StatusError, types.StatusTimeout:
			result.TotalErrors++
		}
	}

	return result
}
