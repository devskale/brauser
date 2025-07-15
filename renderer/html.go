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

	// Extract and print title
	title := doc.Find("title").Text()
	fmt.Println("Title:", title)

	// Extract and print headings
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s: %s\n", s.Get(0).Data, s.Text())
	})

	// Extract and print paragraphs
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	// Extract and print links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			fmt.Printf("Link: %s (%s)\n", s.Text(), href)
		}
	})

	// For HN-specific: Extract story items
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".titleline > a").Text()
		link := s.Find(".titleline > a").AttrOr("href", "")
		fmt.Printf("Story: %s (%s)\n", title, link)
	})

	// Render images as ASCII art
	r.renderImages(doc, baseURL)

	return doc, nil
}

// renderImages processes and renders all images in the document
func (r *HTMLRenderer) renderImages(doc *goquery.Document, baseURL string) {
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			asciiArt, err := r.imageRenderer.RenderImageAsASCII(src, baseURL)
			if err != nil {
				fmt.Printf("Failed to render image %s: %v\n", src, err)
				return
			}
			fmt.Println("ASCII Image:")
			fmt.Println(asciiArt)
		}
	})
}