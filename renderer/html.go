package renderer

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLRenderer handles rendering of HTML content to terminal
type HTMLRenderer struct {
	imageRenderer *ImageRenderer
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		imageRenderer: NewImageRenderer(),
	}
}

// RenderHTML parses the HTML content and displays it in a structured format
func (r *HTMLRenderer) RenderHTML(htmlContent, baseURL string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("           BRAUSER - TERMINAL WEB CONTENT")
	fmt.Println(strings.Repeat("=", 60))

	// Extract and print title
	title := doc.Find("title").Text()
	if title != "" {
		fmt.Printf("\n📄 TITLE: %s\n", title)
		fmt.Println(strings.Repeat("-", len(title)+10))
	}

	// Extract and print headings with hierarchy
	headingCount := 0
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		headingCount++
		tagName := s.Get(0).Data
		text := strings.TrimSpace(s.Text())
		if text != "" {
			switch tagName {
			case "h1":
				fmt.Printf("\n🔸 %s\n", text)
				fmt.Println(strings.Repeat("=", len(text)))
			case "h2":
				fmt.Printf("\n▸ %s\n", text)
				fmt.Println(strings.Repeat("-", len(text)))
			default:
				fmt.Printf("\n• %s\n", text)
			}
		}
	})

	// Extract and print paragraphs with better formatting
	paragraphCount := 0
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 10 { // Filter out very short paragraphs
			paragraphCount++
			fmt.Printf("\n%s\n", text)
		}
	})

	// Extract and print main content areas
	doc.Find("main, article, .content, .main-content, #content").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 50 {
			fmt.Printf("\n📝 MAIN CONTENT:\n%s\n", text[:min(500, len(text))])
			if len(text) > 500 {
				fmt.Println("... (content truncated)")
			}
		}
	})

	// Note: Link extraction and numbering is now handled by the Navigator
	// This renderer focuses on content display only

	// Note: Story extraction is now handled by the Navigator for interactive selection

	// Extract lists
	listCount := 0
	doc.Find("ul li, ol li").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 5 && listCount < 10 {
			listCount++
			if listCount == 1 {
				fmt.Println("\n📋 LIST ITEMS:")
			}
			fmt.Printf("  • %s\n", text)
		}
	})

	// Render images as ASCII art
	r.renderImages(doc, baseURL)

	// Summary
	fmt.Printf("\n" + strings.Repeat("=", 60))
	fmt.Printf("\n📊 CONTENT SUMMARY: %d headings, %d paragraphs\n", headingCount, paragraphCount)
	fmt.Println("💡 Use navigation menu to interact with links")
	fmt.Println(strings.Repeat("=", 60))

	return doc, nil
}

// renderImages processes and renders all images in the document
func (r *HTMLRenderer) renderImages(doc *goquery.Document, baseURL string) {
	imageCount := 0
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		alt := s.AttrOr("alt", "")
		if exists && imageCount < 5 { // Limit to 5 images to avoid spam
			imageCount++
			if imageCount == 1 {
				fmt.Println("\n🖼️  IMAGES:")
			}
			fmt.Printf("  Image %d: %s", imageCount, src)
			if alt != "" {
				fmt.Printf(" (alt: %s)", alt)
			}
			fmt.Println()
			
			// Try to render as ASCII art
			asciiArt, err := r.imageRenderer.RenderImageAsASCII(src, baseURL)
			if err != nil {
				fmt.Printf("    [ASCII conversion failed: %v]\n", err)
			} else {
				fmt.Println("    ASCII Art:")
				fmt.Println(asciiArt)
			}
		}
	})
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}