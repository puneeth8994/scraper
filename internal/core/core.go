package core

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

// Function to scrape JSON
func scrapeJSON(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("non 2xx returned")
		return "", errors.New("Non 2xx response received")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)

	// Access the "title" field and type assert it as a string
	title, ok := result["title"].(string)
	if ok {
		return title, nil
	}

	return "", fmt.Errorf("title field not found or is not a string")
}

// Function to scrape HTML
func scrapeHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("non 2xx returned")
		return "", errors.New("Non 2xx response received")
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	title := fetchTitleFromHTML(doc)

	return title, nil
}

func fetchTitleFromHTML(n *html.Node) string {

	// if we find a h1 tag, we search inside it's attributes.
	// if we find product title class as an attribute, we return the data in the tag
	// Look for <h1> tag with class="product-title"
	if n.Type == html.ElementNode && n.Data == "h1" {

		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "product-title" {
				if n.FirstChild != nil {
					return n.FirstChild.Data
				}
			}
		}
	}

	// In the above condition, we didn't find h1 tag.
	// We start iterating over the child nodes
	// example -
	// <div>
	//	<p>Hello world 1</p>
	//  <p>Hello world 2</p>
	// </div>
	// In the above html, both p tags are siblings
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title := fetchTitleFromHTML(c)
		if title != "" {
			return title
		}
	}

	return ""
}

// Function to read URLs from a file
func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

// ConcurrentScrapingWithRateLimit helps parse http responses concurrently and also applies rate Limiting
func ConcurrentScrapingWithRateLimit(urls []string, rateValue rate.Limit, burstSize int, maxGoroutines int) map[string]string {
	var wg sync.WaitGroup
	limiter := rate.NewLimiter(rateValue, burstSize)

	// Channel for feeding URLs to workers
	urlChannel := make(chan string)
	results := map[string]string{}
	var mu sync.Mutex // To protect the shared results map

	// TODO: Implement batching instead of initiating go routing for every URL
	// Worker function to process URLs from the channel
	worker := func() {
		defer wg.Done()
		for url := range urlChannel {
			// Apply rate limiting
			limiter.Wait(context.Background())

			var title string
			var err error
			if strings.HasSuffix(url, ".json") {
				title, err = scrapeJSON(url)
				if err != nil {
					fmt.Println("Error scraping JSON:", err)
				}
			} else if strings.HasSuffix(url, ".html") {
				title, err = scrapeHTML(url)
				if err != nil {
					fmt.Println("Error scraping HTML:", err)
				}
			} else {
				// TODO: Handle XML formats
			}

			// Store the result safely using mutex
			mu.Lock()
			if err != nil {
				results[url] = "Error: " + err.Error()
			} else {
				results[url] = title
			}
			mu.Unlock()
		}
	}

	// Launch a fixed number of workers (goroutines)
	// to prevent OOM issues
	wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go worker()
	}

	// Send URLs to the channel
	for _, url := range urls {
		urlChannel <- url
	}

	// Close the channel after all URLs have been sent
	close(urlChannel)

	// Wait for all workers to finish
	wg.Wait()

	for key, value := range results {
		// Do something with key and value
		fmt.Println("Key:", key, "Value:", value)
	}

	return results
}
