package renderer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLRenderer handles rendering of HTML content to terminal
type HTMLRenderer struct {
	imageRenderer *ImageRenderer
	outputBuffer  strings.Builder
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		imageRenderer: NewImageRenderer(),
	}
}

// compressEmptyLines removes multiple consecutive empty lines and replaces them with single empty lines
func (r *HTMLRenderer) compressEmptyLines(text string) string {
	// Replace multiple consecutive newlines with double newlines (single empty line)
	re := regexp.MustCompile(`\n\s*\n\s*\n+`)
	return re.ReplaceAllString(text, "\n\n")
}

// printf writes formatted output to the buffer instead of directly to stdout
func (r *HTMLRenderer) printf(format string, args ...interface{}) {
	r.outputBuffer.WriteString(fmt.Sprintf(format, args...))
}

// println writes a line to the buffer instead of directly to stdout
func (r *HTMLRenderer) println(text string) {
	r.outputBuffer.WriteString(text + "\n")
}

// flushOutput compresses empty lines and prints the final output
func (r *HTMLRenderer) flushOutput() {
	compressed := r.compressEmptyLines(r.outputBuffer.String())
	fmt.Print(compressed)
	r.outputBuffer.Reset()
}

// RenderHTML parses the HTML content and displays it in a structured format
func (r *HTMLRenderer) RenderHTML(htmlContent, baseURL string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	r.println("\n" + strings.Repeat("=", 60))
	r.println("           BRAUSER - TERMINAL WEB CONTENT")
	r.println(strings.Repeat("=", 60))

	// Extract and print title
	title := doc.Find("title").Text()
	if title != "" {
		r.printf("\n📄 TITLE: %s\n", title)
		r.println(strings.Repeat("-", len(title)+10))
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
				r.printf("\n🔸 %s\n", text)
				r.println(strings.Repeat("=", len(text)))
			case "h2":
				r.printf("\n▸ %s\n", text)
				r.println(strings.Repeat("-", len(text)))
			default:
				r.printf("\n• %s\n", text)
			}
		}
	})

	// Extract Hacker News stories (table-based layout)
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		if i >= 30 { // Limit to top 30 stories
			return
		}
		headingCount++
		
		// Get story title and link
		titleLink := s.Find(".titleline > a").First()
		title := strings.TrimSpace(titleLink.Text())
		if title != "" {
			r.printf("\n📰 %s\n", title)
		}
	})

	// Extract and print paragraphs with better formatting
	paragraphCount := 0
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 10 { // Filter out very short paragraphs
			paragraphCount++
			r.printf("\n%s\n", text)
		}
	})

	// Extract content from common container elements if no paragraphs found
	if paragraphCount == 0 {
		doc.Find(".content, .main, .post, .entry, .article-content, .story-content").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 50 {
				paragraphCount++
				r.printf("\n📝 CONTENT:\n%s\n", text[:min(300, len(text))])
				if len(text) > 300 {
					r.println("... (content truncated)")
				}
			}
		})
	}

	// Extract and print main content areas
	doc.Find("main, article, .content, .main-content, #content").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 50 {
			r.printf("\n📝 MAIN CONTENT:\n%s\n", text[:min(500, len(text))])
			if len(text) > 500 {
				r.println("... (content truncated)")
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
				r.println("\n📋 LIST ITEMS:")
			}
			r.printf("  • %s\n", text)
		}
	})

	// Render images as ASCII art
	r.renderImages(doc, baseURL)

	// Summary
	r.printf("\n" + strings.Repeat("=", 60))
	r.printf("\n📊 CONTENT SUMMARY: %d headings, %d paragraphs\n", headingCount, paragraphCount)
	r.println("💡 Use navigation menu to interact with links")
	r.println(strings.Repeat("=", 60))

	// Flush the buffered output with compressed empty lines
	r.flushOutput()

	return doc, nil
}

// renderImages processes and renders all images in the document
func (r *HTMLRenderer) renderImages(doc *goquery.Document, baseURL string) {
	imageCount := 0
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		alt := s.AttrOr("alt", "")
		if exists && imageCount < 5 { // Limit to 5 images to avoid spam
			// Skip problematic image formats
			if strings.HasSuffix(strings.ToLower(src), ".svg") ||
			   strings.Contains(strings.ToLower(src), "1x1") ||
			   strings.Contains(strings.ToLower(src), "pixel") {
				return // Skip tracking pixels and SVGs
			}
			
			imageCount++
			if imageCount == 1 {
				r.println("\n🖼️  IMAGES:")
			}
			r.printf("  Image %d: %s", imageCount, src)
			if alt != "" {
				r.printf(" (alt: %s)", alt)
			}
			r.println("")
			
			// Try to render as ASCII art
			asciiArt, err := r.imageRenderer.RenderImageAsASCII(src, baseURL)
			if err != nil {
				if alt != "" {
					r.printf("    [Image: %s]\n", alt)
				} else {
					r.printf("    [Image conversion failed: unsupported format]\n")
				}
			} else {
				r.println("    ASCII Art:")
				r.println(asciiArt)
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