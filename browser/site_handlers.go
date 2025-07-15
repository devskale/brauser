package browser

import (
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// SiteHandler defines the interface for site-specific content handling
type SiteHandler interface {
	CanHandle(url string) bool
	ProcessContent(content string, baseURL string) (string, error)
	GetWaitTime() time.Duration
	RequiresRetry(content string) bool
}

// SiteHandlerManager manages site-specific handlers
type SiteHandlerManager struct {
	handlers []SiteHandler
}

// NewSiteHandlerManager creates a new site handler manager with default handlers
func NewSiteHandlerManager() *SiteHandlerManager {
	manager := &SiteHandlerManager{}
	
	// Register default site handlers
	manager.RegisterHandler(&CodePenHandler{})
	manager.RegisterHandler(&DerStandardHandler{})
	manager.RegisterHandler(&GenericSPAHandler{})
	
	return manager
}

// RegisterHandler adds a new site handler
func (sm *SiteHandlerManager) RegisterHandler(handler SiteHandler) {
	sm.handlers = append(sm.handlers, handler)
}

// GetHandler returns the appropriate handler for a URL
func (sm *SiteHandlerManager) GetHandler(url string) SiteHandler {
	for _, handler := range sm.handlers {
		if handler.CanHandle(url) {
			return handler
		}
	}
	return nil
}

// CodePenHandler handles CodePen-specific content processing
type CodePenHandler struct{}

// CanHandle checks if this handler can process the given URL
func (h *CodePenHandler) CanHandle(url string) bool {
	return strings.Contains(url, "codepen.io")
}

// ProcessContent processes CodePen content to extract meaningful information
func (h *CodePenHandler) ProcessContent(content string, baseURL string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return content, err
	}
	
	// CodePen-specific processing
	// Look for pen title and description
	title := doc.Find("h1.pen-title, .pen-title h1, [data-slug-hash] h1").First().Text()
	description := doc.Find(".pen-description, .pen-details .description").First().Text()
	author := doc.Find(".pen-author, .profile-name").First().Text()
	
	// Extract code if visible
	htmlCode := doc.Find("#html-editor .ace_content, .code-wrap.html-wrap pre").Text()
	cssCode := doc.Find("#css-editor .ace_content, .code-wrap.css-wrap pre").Text()
	jsCode := doc.Find("#js-editor .ace_content, .code-wrap.js-wrap pre").Text()
	
	// Build enhanced content
	var enhanced strings.Builder
	enhanced.WriteString("=== CodePen ===\n")
	
	if title != "" {
		enhanced.WriteString("Title: " + strings.TrimSpace(title) + "\n")
	}
	
	if author != "" {
		enhanced.WriteString("Author: " + strings.TrimSpace(author) + "\n")
	}
	
	if description != "" {
		enhanced.WriteString("Description: " + strings.TrimSpace(description) + "\n")
	}
	
	if htmlCode != "" {
		enhanced.WriteString("\n--- HTML Code ---\n")
		enhanced.WriteString(strings.TrimSpace(htmlCode) + "\n")
	}
	
	if cssCode != "" {
		enhanced.WriteString("\n--- CSS Code ---\n")
		enhanced.WriteString(strings.TrimSpace(cssCode) + "\n")
	}
	
	if jsCode != "" {
		enhanced.WriteString("\n--- JavaScript Code ---\n")
		enhanced.WriteString(strings.TrimSpace(jsCode) + "\n")
	}
	
	// If we found CodePen-specific content, return enhanced version
	if title != "" || htmlCode != "" || cssCode != "" || jsCode != "" {
		return enhanced.String(), nil
	}
	
	// Otherwise return original content
	return content, nil
}

// GetWaitTime returns the recommended wait time for CodePen
func (h *CodePenHandler) GetWaitTime() time.Duration {
	return 3 * time.Second // CodePen often loads content dynamically
}

// RequiresRetry checks if CodePen content needs more time to load
func (h *CodePenHandler) RequiresRetry(content string) bool {
	// Check if we have minimal CodePen content
	return !strings.Contains(content, "pen-title") && !strings.Contains(content, "code-wrap")
}

// DerStandardHandler handles derstandard.at specific content (adblock banners)
type DerStandardHandler struct{}

// CanHandle checks if this handler can process the given URL
func (h *DerStandardHandler) CanHandle(url string) bool {
	return strings.Contains(url, "derstandard.at")
}

// ProcessContent processes DerStandard content to handle adblock banners
func (h *DerStandardHandler) ProcessContent(content string, baseURL string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return content, err
	}
	
	// Check for adblock banner
	adblockBanner := doc.Find(".adblock-banner, #adblock-message, .adblocker-message").Text()
	if adblockBanner != "" {
		// Return a message about the adblock banner
		return "=== DerStandard.at ===\n" +
			"AdBlock Banner Detected: " + strings.TrimSpace(adblockBanner) + "\n" +
			"Note: This site requires disabling ad blockers to view full content.\n", nil
	}
	
	// Extract main article content
	articleTitle := doc.Find("h1, .article-title, .headline").First().Text()
	articleContent := doc.Find(".article-content, .article-body, .content").Text()
	
	if articleTitle != "" {
		var enhanced strings.Builder
		enhanced.WriteString("=== DerStandard.at ===\n")
		enhanced.WriteString("Title: " + strings.TrimSpace(articleTitle) + "\n\n")
		
		if articleContent != "" {
			enhanced.WriteString(strings.TrimSpace(articleContent))
		}
		
		return enhanced.String(), nil
	}
	
	return content, nil
}

// GetWaitTime returns the recommended wait time for DerStandard
func (h *DerStandardHandler) GetWaitTime() time.Duration {
	return 2 * time.Second
}

// RequiresRetry checks if DerStandard content needs more time to load
func (h *DerStandardHandler) RequiresRetry(content string) bool {
	// If we only see adblock banner, might need retry
	return strings.Contains(strings.ToLower(content), "adblock") && len(content) < 1000
}

// GenericSPAHandler handles Single Page Applications that load content dynamically
type GenericSPAHandler struct{}

// CanHandle checks if this handler should process the URL (fallback for SPAs)
func (h *GenericSPAHandler) CanHandle(urlStr string) bool {
	// This is a fallback handler, so it can handle any URL
	// but should be registered last
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	
	// Common SPA indicators in domain names
	spaIndicators := []string{
		"app.",
		"dashboard.",
		"admin.",
		"portal.",
	}
	
	for _, indicator := range spaIndicators {
		if strings.Contains(parsedURL.Host, indicator) {
			return true
		}
	}
	
	return false
}

// ProcessContent processes SPA content
func (h *GenericSPAHandler) ProcessContent(content string, baseURL string) (string, error) {
	// For SPAs, we mainly want to detect if content is loaded
	// and provide better error messages
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return content, err
	}
	
	// Check for common SPA loading indicators
	loadingElements := doc.Find(".loading, .spinner, .loader, [data-loading], #loading")
	if loadingElements.Length() > 0 {
		return "=== Single Page Application ===\n" +
			"Content is still loading. This appears to be a dynamic web application.\n" +
			"Note: Terminal browsers have limited support for dynamic content.\n", nil
	}
	
	// Check for minimal content (likely SPA shell)
	body := doc.Find("body").Text()
	if len(strings.TrimSpace(body)) < 200 {
		return "=== Single Page Application ===\n" +
			"This appears to be a single-page application with minimal initial content.\n" +
			"The main content is likely loaded via JavaScript after page load.\n", nil
	}
	
	return content, nil
}

// GetWaitTime returns the recommended wait time for SPAs
func (h *GenericSPAHandler) GetWaitTime() time.Duration {
	return 5 * time.Second // SPAs often need more time
}

// RequiresRetry checks if SPA content needs more time to load
func (h *GenericSPAHandler) RequiresRetry(content string) bool {
	// Check for minimal content or loading indicators
	lowerContent := strings.ToLower(content)
	return strings.Contains(lowerContent, "loading") ||
		strings.Contains(lowerContent, "spinner") ||
		len(strings.TrimSpace(content)) < 500
}