package main

import (
	"flag"
	"log"
	"net/http"
	"scraper/internal/mock"
	"scraper/internal/services"
	"time"

	"github.com/gorilla/mux"
)

// TODO: Modularize web server as more routes get added
// Function to start the mock web server
func startMockServer() {

	r := mux.NewRouter()
	r.HandleFunc("/entity-{slug}-{uuid}.json", mock.JSONHandler)
	r.HandleFunc("/product-{slug}.html", mock.HTMLHandler)
	r.HandleFunc("/ping", mock.PingHandler)    // Health check URL
	log.Fatal(http.ListenAndServe(":8080", r)) // Start server on port 8080
}

// Function to check server health by making a simple GET request to /ping
func waitForMockServer(url string) {
	for {
		log.Println("Requesting healthcheck to see if server is ready!!!")
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			log.Println("Server is ready!")
			break
		}

		log.Println("Server not yet ready. Waiting!!!")
		time.Sleep(500 * time.Millisecond) // Wait for half a second before retrying
	}
}

func main() {

	// Define flag for input file if applicable
	urlFile := flag.String("file", "", "Specify a file containing URLs")
	flag.Parse()

	// Start mock server in the background
	go startMockServer()

	// waitForMockServer waits for the mockServer by polling health check post which the scraping activity can start
	waitForMockServer("http://localhost:8080/ping")

	services.InitializeScrape(*urlFile)
}
