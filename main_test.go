package main
import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFetchPage tests the fetchPage function with a mock server.
func TestFetchPage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test content"))
	}))
	defer server.Close()

	content, err := fetchPage(server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if content != "Test content" {
		t.Errorf("Expected \"Test content\", got %q", content)
	}
}