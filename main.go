package main

import (
	"fmt"
	"net/http"
)

// CORS middleware to handle CORS requests
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Change "*" to your frontend URL in production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// generateLink handles the form submission
func generateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the link and custom name from the form
		
		customName := r.FormValue("customName")

		// Generate the new link
		newLink := fmt.Sprintf("%s/%s", "bytelink.com/bitly", customName)

		// Prepare the response
		response := fmt.Sprintf("Your new link is: <a href='%s' target='_blank'>%s</a>", newLink, newLink)

		// Send the response back to the client
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, response)
		return
	}

	// If not a POST request, show an error
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func main() {
	// Serve static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Handle the form submission with CORS middleware
	http.Handle("/generate-free-link", cors(http.HandlerFunc(generateLink)))

		// Set Content Security Policy
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' https://unpkg.com; object-src 'none';")
			http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
		})

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil) 
}