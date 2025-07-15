package navigation

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Link represents a clickable link with numbered selection
type Link struct {
	Number int
	Text   string
	URL    string
	Type   string // "nav", "content", "story", etc.
}

// HistoryEntry represents a page in the browser history
type HistoryEntry struct {
	URL     string
	Title   string
	Content string
}

// Navigator handles interactive navigation functionality
type Navigator struct {
	history      []HistoryEntry
	currentIndex int
	links        []Link
	reader       *bufio.Reader
}

// NewNavigator creates a new navigation handler
func NewNavigator() *Navigator {
	return &Navigator{
		history:      make([]HistoryEntry, 0),
		currentIndex: -1,
		links:        make([]Link, 0),
		reader:       bufio.NewReader(os.Stdin),
	}
}

// AddToHistory adds a new page to the browser history
func (n *Navigator) AddToHistory(url, title, content string) {
	// Remove any forward history if we're not at the end
	if n.currentIndex < len(n.history)-1 {
		n.history = n.history[:n.currentIndex+1]
	}
	
	// Add new entry
	entry := HistoryEntry{
		URL:     url,
		Title:   title,
		Content: content,
	}
	n.history = append(n.history, entry)
	n.currentIndex = len(n.history) - 1
	
	// Limit history size to prevent memory issues
	if len(n.history) > 50 {
		n.history = n.history[1:]
		n.currentIndex--
	}
}

// CanGoBack returns true if there's a previous page in history
func (n *Navigator) CanGoBack() bool {
	return n.currentIndex > 0
}

// CanGoForward returns true if there's a next page in history
func (n *Navigator) CanGoForward() bool {
	return n.currentIndex < len(n.history)-1
}

// GoBack navigates to the previous page in history
func (n *Navigator) GoBack() *HistoryEntry {
	if !n.CanGoBack() {
		return nil
	}
	n.currentIndex--
	return &n.history[n.currentIndex]
}

// GoForward navigates to the next page in history
func (n *Navigator) GoForward() *HistoryEntry {
	if !n.CanGoForward() {
		return nil
	}
	n.currentIndex++
	return &n.history[n.currentIndex]
}

// GetCurrentPage returns the current page from history
func (n *Navigator) GetCurrentPage() *HistoryEntry {
	if n.currentIndex >= 0 && n.currentIndex < len(n.history) {
		return &n.history[n.currentIndex]
	}
	return nil
}

// cleanLinkText removes excessive whitespace and newlines from link text
func (n *Navigator) cleanLinkText(text string) string {
	// Remove all newlines and replace with single spaces
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(text, " ")
	return strings.TrimSpace(cleaned)
}

// ExtractLinks extracts all clickable links from the HTML document and assigns numbers
func (n *Navigator) ExtractLinks(doc *goquery.Document, baseURL string) {
	n.links = make([]Link, 0)
	linkNumber := 1
	
	// Parse base URL for resolving relative links
	base, err := url.Parse(baseURL)
	if err != nil {
		base = nil
	}
	
	// Extract navigation links
	doc.Find("nav a, .nav a, .menu a").Each(func(i int, s *goquery.Selection) {
		text := n.cleanLinkText(s.Text())
		href, exists := s.Attr("href")
		if exists && text != "" {
			resolvedURL := n.resolveURL(href, base)
			n.links = append(n.links, Link{
				Number: linkNumber,
				Text:   text,
				URL:    resolvedURL,
				Type:   "nav",
			})
			linkNumber++
		}
	})
	
	// Extract content links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := n.cleanLinkText(s.Text())
		href, exists := s.Attr("href")
		if exists && text != "" && len(text) > 3 && linkNumber <= 50 {
			// Skip if already added as navigation
			alreadyAdded := false
			for _, link := range n.links {
				if link.URL == href {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				resolvedURL := n.resolveURL(href, base)
				n.links = append(n.links, Link{
					Number: linkNumber,
					Text:   text,
					URL:    resolvedURL,
					Type:   "content",
				})
				linkNumber++
			}
		}
	})
	
	// Extract HackerNews stories
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		title := n.cleanLinkText(s.Find(".titleline > a").Text())
		link := s.Find(".titleline > a").AttrOr("href", "")
		if title != "" && link != "" && linkNumber <= 50 {
			resolvedURL := n.resolveURL(link, base)
			n.links = append(n.links, Link{
				Number: linkNumber,
				Text:   title,
				URL:    resolvedURL,
				Type:   "story",
			})
			linkNumber++
		}
	})
}

// resolveURL resolves relative URLs against the base URL
func (n *Navigator) resolveURL(href string, base *url.URL) string {
	if base == nil {
		return href
	}
	
	parsed, err := url.Parse(href)
	if err != nil {
		return href
	}
	
	return base.ResolveReference(parsed).String()
}

// DisplayLinks shows all numbered links to the user
func (n *Navigator) DisplayLinks() {
	if len(n.links) == 0 {
		fmt.Println("\n❌ No clickable links found on this page.")
		return
	}
	
	fmt.Printf("\n🔗 CLICKABLE LINKS (%d total):\n", len(n.links))
	fmt.Println(strings.Repeat("-", 50))
	
	// Group links by type
	navLinks := make([]Link, 0)
	contentLinks := make([]Link, 0)
	storyLinks := make([]Link, 0)
	
	for _, link := range n.links {
		switch link.Type {
		case "nav":
			navLinks = append(navLinks, link)
		case "story":
			storyLinks = append(storyLinks, link)
		default:
			contentLinks = append(contentLinks, link)
		}
	}
	
	// Display navigation links
	if len(navLinks) > 0 {
		fmt.Println("\n🧭 Navigation:")
		for _, link := range navLinks {
			fmt.Printf("  [%d] %s\n", link.Number, link.Text)
		}
	}
	
	// Display story links (for sites like HN)
	if len(storyLinks) > 0 {
		fmt.Println("\n📰 Stories:")
		for _, link := range storyLinks {
			fmt.Printf("  [%d] %s\n", link.Number, link.Text)
		}
	}
	
	// Display content links
	if len(contentLinks) > 0 {
		fmt.Println("\n📄 Content Links:")
		for _, link := range contentLinks {
			fmt.Printf("  [%d] %s\n", link.Number, link.Text)
		}
	}
}

// GetLinkByNumber returns the link with the specified number
func (n *Navigator) GetLinkByNumber(number int) *Link {
	for _, link := range n.links {
		if link.Number == number {
			return &link
		}
	}
	return nil
}

// GetLinks returns all extracted links (for testing purposes)
func (n *Navigator) GetLinks() []Link {
	return n.links
}

// ShowNavigationMenu displays the interactive navigation menu
func (n *Navigator) ShowNavigationMenu() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("           BRAUSER NAVIGATION MENU")
	fmt.Println(strings.Repeat("=", 60))
	
	// Show current page info
	current := n.GetCurrentPage()
	if current != nil {
		fmt.Printf("📍 Current: %s\n", current.URL)
		if current.Title != "" {
			fmt.Printf("📄 Title: %s\n", current.Title)
		}
	}
	
	// Show navigation options
	fmt.Println("\n🎯 Navigation Options:")
	fmt.Println("  • Type a number [1-50] to follow a link")
	fmt.Println("  • Type 'b' or 'back' to go back")
	fmt.Println("  • Type 'f' or 'forward' to go forward")
	fmt.Println("  • Type 'h' or 'history' to view history")
	fmt.Println("  • Type 'l' or 'links' to show links again")
	fmt.Println("  • Type 'u' or 'url' to enter a new URL")
	fmt.Println("  • Type 'r' or 'refresh' to reload current page")
	fmt.Println("  • Type 'q' or 'quit' to exit")
	
	// Show back/forward status
	if n.CanGoBack() {
		fmt.Println("  ⬅️  Back available")
	}
	if n.CanGoForward() {
		fmt.Println("  ➡️  Forward available")
	}
	
	fmt.Println(strings.Repeat("-", 60))
}

// GetUserInput prompts the user for input and returns the command
func (n *Navigator) GetUserInput() (string, error) {
	fmt.Print("\n🌐 brauser> ")
	input, err := n.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ProcessUserInput processes user input and returns the action to take
func (n *Navigator) ProcessUserInput(input string) (action string, data interface{}) {
	input = strings.ToLower(strings.TrimSpace(input))
	
	// Handle numeric input (link selection)
	if num, err := strconv.Atoi(input); err == nil {
		if link := n.GetLinkByNumber(num); link != nil {
			return "navigate", link.URL
		} else {
			return "error", fmt.Sprintf("Link number %d not found. Please choose a number between 1 and %d.", num, len(n.links))
		}
	}
	
	// Handle text commands
	switch input {
	case "b", "back":
		if n.CanGoBack() {
			return "back", nil
		} else {
			return "error", "No previous page in history."
		}
	case "f", "forward":
		if n.CanGoForward() {
			return "forward", nil
		} else {
			return "error", "No next page in history."
		}
	case "h", "history":
		return "history", nil
	case "l", "links":
		return "links", nil
	case "u", "url":
		return "url", nil
	case "r", "refresh":
		return "refresh", nil
	case "q", "quit":
		return "quit", nil
	default:
		return "error", fmt.Sprintf("Unknown command: %s. Type 'h' for help.", input)
	}
}

// ShowHistory displays the browser history
func (n *Navigator) ShowHistory() {
	if len(n.history) == 0 {
		fmt.Println("\n📚 History is empty.")
		return
	}
	
	fmt.Printf("\n📚 BROWSER HISTORY (%d pages):\n", len(n.history))
	fmt.Println(strings.Repeat("-", 50))
	
	for i, entry := range n.history {
		marker := "  "
		if i == n.currentIndex {
			marker = "👉"
		}
		
		title := entry.Title
		if title == "" {
			title = "(No title)"
		}
		
		fmt.Printf("%s %d. %s\n", marker, i+1, title)
		fmt.Printf("     %s\n", entry.URL)
	}
}

// PromptForURL prompts the user to enter a new URL
func (n *Navigator) PromptForURL() (string, error) {
	fmt.Print("\n🌐 Enter URL: ")
	input, err := n.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	url := strings.TrimSpace(input)
	if url == "" {
		return "", fmt.Errorf("empty URL")
	}
	
	// Add http:// if no protocol specified
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	
	return url, nil
}