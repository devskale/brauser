package main

import (
	"fmt"
	"log"
	"os"

	"brauser/browser"
	"brauser/js"
	"brauser/renderer"
)

// main is the entry point of the Brauser application.
func main() {
	fmt.Println("Brauser: Minimalistic Terminal Web Browser")
	if len(os.Args) < 2 {
		fmt.Println("Usage: brauser <url> [--no-retry]")
		fmt.Println("  --no-retry: Disable content detection and retry logic")
		return
	}
	
	url := os.Args[1]
	enableRetry := true
	
	// Check for --no-retry flag
	if len(os.Args) > 2 && os.Args[2] == "--no-retry" {
		enableRetry = false
		fmt.Println("Content detection and retry logic disabled")
	}
	
	// Create browser client
	client := browser.NewClient()
	
	// Fetch page content with enhanced detection
	var content string
	var analysis *browser.ContentAnalysis
	var err error
	
	if enableRetry {
		content, analysis, err = client.FetchPageWithEnhancedDetection(url)
		if err != nil {
			log.Fatalf("Failed to fetch page: %v", err)
		}
		
		// Display content analysis results
		displayContentAnalysis(analysis)
	} else {
		content, err = client.FetchPageWithRetry(url, false)
		if err != nil {
			log.Fatalf("Failed to fetch page: %v", err)
		}
	}
	
	// Create HTML renderer
	htmlRenderer := renderer.NewHTMLRenderer()
	
	// Render HTML content
	doc, err := htmlRenderer.RenderHTML(content, url)
	if err != nil {
		log.Fatalf("Failed to render HTML: %v", err)
	}
	
	// Execute embedded JavaScript
	title := doc.Find("title").Text()
	js.ExecuteJS(doc, title)
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