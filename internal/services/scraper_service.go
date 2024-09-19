package services

import (
	"log"
	"scraper/internal/core"
	"scraper/internal/utils"
)

// Default URLs to scrape if no urlFile is provided
var defaultScrapeURL = []string{
	"http://localhost:8080/entity-book-f3bfa24c-2645-48c0-9117-b338bef9b9ab.json",
	"http://localhost:8080/product-book.html",
}

// InitializeScrape initializes scrape and access url file as an argument
// if the urlFile is not present then it uses the defaultScrapeURL to scrape
func InitializeScrape(urlFile string) {

	// If urlfile is a non empty string, we will try to read it
	if urlFile != "" {
		urls, err := utils.ReadURLsFromFile(urlFile)
		if err != nil {
			log.Fatalf("Error reading URLs from file: %v", err)
		}
		// Limit: 1 request per second, burst size of 2
		core.ConcurrentScrapingWithRateLimit(urls, 1, 2)
	} else {
		// Limit: 1 request per second, burst size of 2
		core.ConcurrentScrapingWithRateLimit(defaultScrapeURL, 1, 2)
	}
}
