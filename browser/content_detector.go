package browser

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ContentDetector analyzes page content to determine if it's fully loaded
type ContentDetector struct {
	loadingIndicators []string
	loadingPatterns   []*regexp.Regexp
	minContentLength  int
}

// NewContentDetector creates a new content detector with default patterns
func NewContentDetector() *ContentDetector {
	return &ContentDetector{
		loadingIndicators: []string{
			"just a moment",
			"loading",
			"please wait",
			"redirecting",
			"checking your browser",
			"verifying you are human",
			"one moment please",
			"loading content",
			"initializing",
			"preparing",
		},
		loadingPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)loading[\s\.]{0,3}`),
			regexp.MustCompile(`(?i)please\s+wait`),
			regexp.MustCompile(`(?i)just\s+a\s+moment`),
			regexp.MustCompile(`(?i)checking\s+your\s+browser`),
			regexp.MustCompile(`(?i)verifying\s+you\s+are\s+human`),
			regexp.MustCompile(`(?i)cloudflare`),
			regexp.MustCompile(`(?i)ddos\s+protection`),
			regexp.MustCompile(`(?i)security\s+check`),
		},
		minContentLength: 500, // Minimum content length to consider page loaded
	}
}

// ContentAnalysis represents the result of content analysis
type ContentAnalysis struct {
	IsLoaded           bool
	IsLoadingPage      bool
	IsCookieBanner     bool
	IsAdBlockBanner    bool
	IsInterstitial     bool
	ContentLength      int
	LoadingIndicators  []string
	SuggestedWaitTime  time.Duration
	RequiresRetry      bool
}

// AnalyzeContent analyzes the HTML content to determine its state
func (cd *ContentDetector) AnalyzeContent(htmlContent string) *ContentAnalysis {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return &ContentAnalysis{
			IsLoaded:      false,
			ContentLength: len(htmlContent),
		}
	}

	analysis := &ContentAnalysis{
		ContentLength: len(htmlContent),
	}

	// Extract visible text content
	visibleText := cd.extractVisibleText(doc)
	analysis.ContentLength = len(visibleText)

	// Check for loading indicators
	analysis.LoadingIndicators = cd.findLoadingIndicators(visibleText)
	analysis.IsLoadingPage = len(analysis.LoadingIndicators) > 0

	// Check for cookie banners
	analysis.IsCookieBanner = cd.detectCookieBanner(doc, visibleText)

	// Check for adblock banners
	analysis.IsAdBlockBanner = cd.detectAdBlockBanner(doc, visibleText)

	// Check for interstitial pages
	analysis.IsInterstitial = analysis.IsCookieBanner || analysis.IsAdBlockBanner || analysis.IsLoadingPage

	// Determine if content is loaded
	analysis.IsLoaded = cd.isContentLoaded(analysis)

	// Suggest wait time if needed
	analysis.SuggestedWaitTime = cd.calculateWaitTime(analysis)

	// Determine if retry is needed
	analysis.RequiresRetry = !analysis.IsLoaded && (analysis.IsLoadingPage || analysis.ContentLength < cd.minContentLength)

	return analysis
}

// extractVisibleText extracts visible text content from the document
func (cd *ContentDetector) extractVisibleText(doc *goquery.Document) string {
	// Remove script and style elements
	doc.Find("script, style, noscript").Remove()
	
	// Extract text from body
	body := doc.Find("body")
	if body.Length() == 0 {
		// Fallback to entire document if no body
		return doc.Text()
	}
	
	return strings.TrimSpace(body.Text())
}

// findLoadingIndicators searches for loading indicators in the text
func (cd *ContentDetector) findLoadingIndicators(text string) []string {
	var indicators []string
	lowerText := strings.ToLower(text)
	
	// Check simple string indicators
	for _, indicator := range cd.loadingIndicators {
		if strings.Contains(lowerText, indicator) {
			indicators = append(indicators, indicator)
		}
	}
	
	// Check regex patterns
	for _, pattern := range cd.loadingPatterns {
		if pattern.MatchString(text) {
			match := pattern.FindString(text)
			indicators = append(indicators, match)
		}
	}
	
	return indicators
}

// detectCookieBanner checks for cookie consent banners
func (cd *ContentDetector) detectCookieBanner(doc *goquery.Document, text string) bool {
	lowerText := strings.ToLower(text)
	
	// Common cookie banner text patterns
	cookiePatterns := []string{
		"accept cookies",
		"cookie policy",
		"we use cookies",
		"cookies help us",
		"cookie consent",
		"privacy policy",
		"accept all",
		"manage cookies",
	}
	
	for _, pattern := range cookiePatterns {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}
	
	// Check for common cookie banner CSS classes/IDs
	cookieSelectors := []string{
		"#cookie-banner",
		".cookie-banner",
		"#cookie-consent",
		".cookie-consent",
		"#gdpr-banner",
		".gdpr-banner",
		"[data-cookie]",
	}
	
	for _, selector := range cookieSelectors {
		if doc.Find(selector).Length() > 0 {
			return true
		}
	}
	
	return false
}

// detectAdBlockBanner checks for adblock detection banners
func (cd *ContentDetector) detectAdBlockBanner(doc *goquery.Document, text string) bool {
	lowerText := strings.ToLower(text)
	
	// Common adblock banner text patterns
	adblockPatterns := []string{
		"disable adblock",
		"turn off adblock",
		"ad blocker detected",
		"please disable",
		"whitelist this site",
		"support us by disabling",
		"ads help us",
	}
	
	for _, pattern := range adblockPatterns {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}
	
	return false
}

// isContentLoaded determines if the page content is fully loaded
func (cd *ContentDetector) isContentLoaded(analysis *ContentAnalysis) bool {
	// If there are loading indicators, content is not loaded
	if analysis.IsLoadingPage {
		return false
	}
	
	// If content is too short, it might be a loading page
	if analysis.ContentLength < cd.minContentLength {
		return false
	}
	
	// If it's an interstitial page, content is not the main content
	if analysis.IsInterstitial {
		return false
	}
	
	return true
}

// calculateWaitTime suggests how long to wait before retrying
func (cd *ContentDetector) calculateWaitTime(analysis *ContentAnalysis) time.Duration {
	if !analysis.RequiresRetry {
		return 0
	}
	
	// Base wait time
	waitTime := 2 * time.Second
	
	// Increase wait time for specific scenarios
	if analysis.IsLoadingPage {
		waitTime = 3 * time.Second
	}
	
	if analysis.IsCookieBanner || analysis.IsAdBlockBanner {
		waitTime = 1 * time.Second // These usually resolve quickly
	}
	
	// Very short content might need more time
	if analysis.ContentLength < 100 {
		waitTime = 5 * time.Second
	}
	
	return waitTime
}

// SetMinContentLength sets the minimum content length threshold
func (cd *ContentDetector) SetMinContentLength(length int) {
	cd.minContentLength = length
}

// AddLoadingIndicator adds a custom loading indicator pattern
func (cd *ContentDetector) AddLoadingIndicator(indicator string) {
	cd.loadingIndicators = append(cd.loadingIndicators, strings.ToLower(indicator))
}

// AddLoadingPattern adds a custom loading regex pattern
func (cd *ContentDetector) AddLoadingPattern(pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	cd.loadingPatterns = append(cd.loadingPatterns, regex)
	return nil
}