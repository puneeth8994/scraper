package mock

import (
	"net/http"

	"github.com/gorilla/mux"
)

// JSONHandler sends mock Mock JSON handler
func JSONHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // Extract variables from the route
	slug := vars["slug"]
	uuid := vars["uuid"]

	_ = slug
	_ = uuid
	// Use these values for internal logic if needed, but don't write them directly to the response
	// fmt.Printf("Processing request for Slug: %s, UUID: %s\n", slug, uuid)

	// Now send only the JSON response
	jsonResponse := `{"title":"My title"}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}

// HTMLHandler mocks a http response
func HTMLHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
		<html>
			<body>
				<h1 class="product-title" data-id="f3bfa24c-2645-48c0-9117-b338bef9b9ab">Product title</h1>
			</body>
		</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlContent))
}

// PingHandler is dedicated handler for health check
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong")) // Simple response
}
