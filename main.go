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
		fmt.Println("Usage: brauser <url>")
		return
	}
	url := os.Args[1]
	
	// Create browser client
	client := browser.NewClient()
	
	// Fetch page content
	content, err := client.FetchPage(url)
	if err != nil {
		log.Fatalf("Failed to fetch page: %v", err)
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