package browser

import (
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client for fetching web pages
type Client struct {
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new browser client with default settings
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
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

// FetchPage fetches the content of the given URL and returns it as a string
func (c *Client) FetchPage(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}