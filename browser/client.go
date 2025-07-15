package browser

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client represents an HTTP client for fetching web pages
type Client struct {
	httpClient      *http.Client
	userAgent       string
	contentDetector *ContentDetector
	siteHandlers    *SiteHandlerManager
	maxRetries      int
	maxWaitTime     time.Duration
}

// NewClient creates a new browser client with default settings
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		contentDetector: NewContentDetector(),
		siteHandlers:    NewSiteHandlerManager(),
		maxRetries:      3,
		maxWaitTime:     10 * time.Second,
	}
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// SetUserAgent sets the User-Agent header
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// SetMaxRetries sets the maximum number of retry attempts
func (c *Client) SetMaxRetries(retries int) {
	c.maxRetries = retries
}

// SetMaxWaitTime sets the maximum wait time for dynamic content
func (c *Client) SetMaxWaitTime(waitTime time.Duration) {
	c.maxWaitTime = waitTime
}

// FetchPage fetches the content of the given URL and returns it as a string
func (c *Client) FetchPage(url string) (string, error) {
	return c.FetchPageWithRetry(url, true)
}

// FetchPageWithRetry fetches content with optional retry logic for dynamic content
func (c *Client) FetchPageWithRetry(url string, enableRetry bool) (string, error) {
	var lastContent string
	
	// Check for site-specific handler
	siteHandler := c.siteHandlers.GetHandler(url)
	
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Fetch the page content
		content, err := c.fetchPageOnce(url)
		if err != nil {
			return "", err
		}
		
		lastContent = content
		
		// Apply site-specific processing if available
		if siteHandler != nil {
			processedContent, err := siteHandler.ProcessContent(content, url)
			if err == nil {
				content = processedContent
				lastContent = content
			}
		}
		
		// Analyze content if retry is enabled
		if !enableRetry {
			return content, nil
		}
		
		analysis := c.contentDetector.AnalyzeContent(content)
		
		// Log analysis results for debugging
		if attempt == 0 {
			c.logContentAnalysis(url, analysis)
		}
		
		// Check site-specific retry logic
		siteNeedsRetry := false
		if siteHandler != nil {
			siteNeedsRetry = siteHandler.RequiresRetry(content)
		}
		
		// If content is loaded or we've reached max retries, return
		if (analysis.IsLoaded && !siteNeedsRetry) || attempt == c.maxRetries {
			return content, nil
		}
		
		// Wait before retrying if content needs more time
		if analysis.RequiresRetry || siteNeedsRetry {
			waitTime := analysis.SuggestedWaitTime
			
			// Use site-specific wait time if available
			if siteHandler != nil {
				siteWaitTime := siteHandler.GetWaitTime()
				if siteWaitTime > waitTime {
					waitTime = siteWaitTime
				}
			}
			
			if waitTime > c.maxWaitTime {
				waitTime = c.maxWaitTime
			}
			
			log.Printf("Content not fully loaded, waiting %v before retry %d/%d", waitTime, attempt+1, c.maxRetries)
			time.Sleep(waitTime)
		} else {
			// No retry needed, return current content
			return content, nil
		}
	}
	
	return lastContent, nil
}

// fetchPageOnce performs a single HTTP request to fetch page content
func (c *Client) fetchPageOnce(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	// Handle gzip decompression
	var reader io.Reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}
	
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

// logContentAnalysis logs the results of content analysis for debugging
func (c *Client) logContentAnalysis(url string, analysis *ContentAnalysis) {
	log.Printf("Content Analysis for %s:", url)
	log.Printf("  - Content Length: %d", analysis.ContentLength)
	log.Printf("  - Is Loaded: %t", analysis.IsLoaded)
	log.Printf("  - Is Loading Page: %t", analysis.IsLoadingPage)
	log.Printf("  - Is Cookie Banner: %t", analysis.IsCookieBanner)
	log.Printf("  - Is AdBlock Banner: %t", analysis.IsAdBlockBanner)
	log.Printf("  - Is Interstitial: %t", analysis.IsInterstitial)
	log.Printf("  - Requires Retry: %t", analysis.RequiresRetry)
	
	if len(analysis.LoadingIndicators) > 0 {
		log.Printf("  - Loading Indicators: %v", analysis.LoadingIndicators)
	}
	
	if analysis.SuggestedWaitTime > 0 {
		log.Printf("  - Suggested Wait Time: %v", analysis.SuggestedWaitTime)
	}
}

// GetContentDetector returns the content detector for customization
func (c *Client) GetContentDetector() *ContentDetector {
	return c.contentDetector
}

// GetSiteHandlers returns the site handler manager for customization
func (c *Client) GetSiteHandlers() *SiteHandlerManager {
	return c.siteHandlers
}

// FetchPageWithEnhancedDetection fetches content with full content detection and site-specific handling
func (c *Client) FetchPageWithEnhancedDetection(url string) (string, *ContentAnalysis, error) {
	content, err := c.FetchPageWithRetry(url, true)
	if err != nil {
		return "", nil, err
	}
	
	analysis := c.contentDetector.AnalyzeContent(content)
	return content, analysis, nil
}