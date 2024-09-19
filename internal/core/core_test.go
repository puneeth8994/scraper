package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScrapeJSON_Success(t *testing.T) {
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

	assert.Equal(t, "Test Title", title)

	if title != "Test Title" {
		t.Errorf("Expected 'Test Title', got %s", title)
	}
}

func TestScrapeJSON_Failure(t *testing.T) {
	// Create a test server that responds with a mock JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"title": "Test Title"}`))
	}))
	defer ts.Close()

	// Test scraping the mock server
	title, err := scrapeJSON(ts.URL)

	assert.NotEmpty(t, err)

	assert.Empty(t, title)
}

func TestScrapeHTML_Success(t *testing.T) {
	// Create a test server that responds with a mock HTML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head></head><body><h1 class="product-title">Product Title</h1></body></html>`))
	}))
	defer ts.Close()

	// Test scraping the mock server
	title, err := scrapeHTML(ts.URL)

	assert.Nil(t, err)

	assert.Equal(t, "Product Title", title)
}

func TestScrapeHTML_Failure(t *testing.T) {
	// Create a test server that responds with a mock HTML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head></head><body><h1 class="product-title">Product Title</h1></body></html>`))
	}))
	defer ts.Close()

	// Test scraping the mock server
	title, err := scrapeHTML(ts.URL)

	assert.NotNil(t, err)

	assert.Empty(t, "", title)
}

func TestConcurrentScrapingWithRateLimit(t *testing.T) {
	// Create a test server that handles both JSON and HTML requests
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".json") {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"title": "Test JSON Title"}`))
		} else if strings.HasSuffix(r.URL.Path, ".html") {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<html><head></head><body><h1 class="product-title">Test HTML Title</h1></body></html>`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	// Mock URLs (some ending in .json and some in .html)
	urls := []string{
		ts.URL + "/test1.json",
		ts.URL + "/test2.html",
		ts.URL + "/test3.json",
		ts.URL + "/test4.html",
	}

	// Test Case 1: Rate Limit = 1 request/sec, Burst Size = 1
	startTime := time.Now()
	ConcurrentScrapingWithRateLimit(urls, 1, 1)
	duration := time.Since(startTime)

	// Ensure that rate limiting was applied correctly (should take ~4 seconds for 4 URLs)
	if duration < 3*time.Second || duration >= 5*time.Second {
		t.Errorf("Rate limiting not applied correctly, execution time too short/long: %v", duration)
	}

	// Test Case 2: Rate Limit = 1 request/sec, Burst Size = 2
	startTime = time.Now()
	ConcurrentScrapingWithRateLimit(urls, 1, 2)
	duration = time.Since(startTime)

	// Ensure that rate limiting was applied correctly (should take ~2 seconds due to burst size)
	if duration < 2*time.Second || duration >= 4*time.Second {
		t.Errorf("Rate limiting with burst size 2 not applied correctly, execution time too short/long: %v", duration)
	}
}
