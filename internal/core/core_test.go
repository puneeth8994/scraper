package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrapeJSON(t *testing.T) {
	// Create a test server that responds with a mock JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"title": "Test Title"}`))
	}))
	defer ts.Close()

	// Test scraping the mock server
	title, err := scrapeJSON(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if title != "Test Title" {
		t.Errorf("Expected 'Test Title', got %s", title)
	}
}
