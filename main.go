package main

import (
	"fmt"
	"os"
	"strings"

	"brauser/browser"
	"brauser/js"
	"brauser/navigation"
	"brauser/renderer"
	"github.com/PuerkitoBio/goquery"
)

// main is the entry point of the Brauser application with interactive navigation.
func main() {
	fmt.Println("Brauser: Minimalistic Terminal Web Browser with Interactive Navigation")
	if len(os.Args) < 2 {
		fmt.Println("Usage: brauser <url> [--no-retry]")
		fmt.Println("  --no-retry: Disable content detection and retry logic")
		fmt.Println("  Interactive features: numbered links, back/forward, URL bar")
		return
	}
	
	initialURL := os.Args[1]
	enableRetry := true
	
	// Check for --no-retry flag
	if len(os.Args) > 2 && os.Args[2] == "--no-retry" {
		enableRetry = false
		fmt.Println("Content detection and retry logic disabled")
	}
	
	// Create components
	client := browser.NewClient()
	htmlRenderer := renderer.NewHTMLRenderer()
	navigator := navigation.NewNavigator()
	
	// Start interactive browsing session
	startInteractiveBrowsing(client, htmlRenderer, navigator, initialURL, enableRetry)
}

// startInteractiveBrowsing handles the main interactive browsing loop
func startInteractiveBrowsing(client *browser.Client, htmlRenderer *renderer.HTMLRenderer, navigator *navigation.Navigator, initialURL string, enableRetry bool) {
	currentURL := initialURL
	
	for {
		// Fetch and display page
		if err := loadAndDisplayPage(client, htmlRenderer, navigator, currentURL, enableRetry); err != nil {
			fmt.Printf("❌ Error loading page: %v\n", err)
			continue
		}
		
		// Show navigation menu and get user input
		navigator.ShowNavigationMenu()
		navigator.DisplayLinks()
		
		for {
			input, err := navigator.GetUserInput()
			if err != nil {
				fmt.Printf("❌ Error reading input: %v\n", err)
				continue
			}
			
			action, data := navigator.ProcessUserInput(input)
			
			switch action {
			case "navigate":
				currentURL = data.(string)
				fmt.Printf("🌐 Navigating to: %s\n", currentURL)
				goto loadPage // Break inner loop and load new page
				
			case "back":
				if entry := navigator.GoBack(); entry != nil {
					currentURL = entry.URL
					fmt.Printf("⬅️  Going back to: %s\n", currentURL)
					// Display cached content and re-extract links
					displayCachedPage(entry, htmlRenderer, navigator)
					navigator.ShowNavigationMenu()
					navigator.DisplayLinks()
				}
				
			case "forward":
				if entry := navigator.GoForward(); entry != nil {
					currentURL = entry.URL
					fmt.Printf("➡️  Going forward to: %s\n", currentURL)
					// Display cached content and re-extract links
					displayCachedPage(entry, htmlRenderer, navigator)
					navigator.ShowNavigationMenu()
					navigator.DisplayLinks()
				}
				
			case "history":
				navigator.ShowHistory()
				
			case "links":
				navigator.DisplayLinks()
				
			case "url":
				newURL, err := navigator.PromptForURL()
				if err != nil {
					fmt.Printf("❌ Error: %v\n", err)
				} else {
					currentURL = newURL
					fmt.Printf("🌐 Navigating to: %s\n", currentURL)
					goto loadPage
				}
				
			case "refresh":
				fmt.Printf("🔄 Refreshing: %s\n", currentURL)
				goto loadPage
				
			case "quit":
				fmt.Println("👋 Thanks for using Brauser!")
				return
				
			case "error":
				fmt.Printf("❌ %s\n", data.(string))
			}
		}
		
		loadPage:
		// Continue to next iteration to load the new page
	}
}

// loadAndDisplayPage fetches, renders, and processes a web page
func loadAndDisplayPage(client *browser.Client, htmlRenderer *renderer.HTMLRenderer, navigator *navigation.Navigator, url string, enableRetry bool) error {
	// Fetch page content
	var content string
	var analysis *browser.ContentAnalysis
	var err error
	
	if enableRetry {
		content, analysis, err = client.FetchPageWithEnhancedDetection(url)
		if err != nil {
			return fmt.Errorf("failed to fetch page: %v", err)
		}
		
		// Display content analysis results
		displayContentAnalysis(analysis)
	} else {
		content, err = client.FetchPageWithRetry(url, false)
		if err != nil {
			return fmt.Errorf("failed to fetch page: %v", err)
		}
	}
	
	// Render HTML content
	doc, err := htmlRenderer.RenderHTML(content, url)
	if err != nil {
		return fmt.Errorf("failed to render HTML: %v", err)
	}
	
	// Execute embedded JavaScript
	title := doc.Find("title").Text()
	js.ExecuteJS(doc, title)
	
	// Extract links for navigation
	navigator.ExtractLinks(doc, url)
	
	// Add to history
	navigator.AddToHistory(url, title, content)
	
	return nil
}

// displayCachedPage shows a cached page from history and re-extracts links
func displayCachedPage(entry *navigation.HistoryEntry, htmlRenderer *renderer.HTMLRenderer, navigator *navigation.Navigator) {
	// Re-render the cached HTML content to display it properly
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(entry.Content))
	if err != nil {
		fmt.Printf("❌ Error parsing cached content: %v\n", err)
		return
	}
	
	// Render the cached content
	_, err = htmlRenderer.RenderHTML(entry.Content, entry.URL)
	if err != nil {
		fmt.Printf("❌ Error rendering cached content: %v\n", err)
		return
	}
	
	// Re-extract links from the cached content
	navigator.ExtractLinks(doc, entry.URL)
	
	fmt.Println("\n💾 (Displaying cached content - use 'r' to refresh)")
}

// displayContentAnalysis shows the content analysis results to the user
func displayContentAnalysis(analysis *browser.ContentAnalysis) {
	if analysis == nil {
		return
	}
	
	fmt.Println("\n=== Content Analysis ===")
	
	if analysis.IsLoaded {
		fmt.Println("✅ Content appears to be fully loaded")
	} else {
		fmt.Println("⚠️  Content may not be fully loaded")
	}
	
	if analysis.IsLoadingPage {
		fmt.Println("🔄 Loading page detected")
		if len(analysis.LoadingIndicators) > 0 {
			fmt.Printf("   Indicators: %v\n", analysis.LoadingIndicators)
		}
	}
	
	if analysis.IsCookieBanner {
		fmt.Println("🍪 Cookie consent banner detected")
	}
	
	if analysis.IsAdBlockBanner {
		fmt.Println("🚫 AdBlock detection banner found")
	}
	
	if analysis.IsInterstitial {
		fmt.Println("📄 Interstitial page detected (banner/loading screen)")
	}
	
	fmt.Printf("📊 Content length: %d characters\n", analysis.ContentLength)
	
	if analysis.RequiresRetry {
		fmt.Printf("🔁 Retry recommended (wait time: %v)\n", analysis.SuggestedWaitTime)
	}
	
	fmt.Println("========================\n")
}