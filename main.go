package main

import (
    "fmt"
    "net/http"
    "os"
    "sync"
)

// Global map to store custom names and their corresponding links
var linkMap = make(map[string]string) // Global map to store custom names and their corresponding links
var mu sync.Mutex // Mutex to handle concurrent access

// serveHTML serves the HTML content from the index.html file
func serveHTML(w http.ResponseWriter, r *http.Request) {
    // Serve the HTML file
    http.ServeFile(w, r, "index.html")
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

        // Get the custom name from the form
        customName := r.FormValue("customName")
        link := r.FormValue("link") // Get the actual link from the form

        // Store the mapping
        mu.Lock()
        linkMap[customName] = link
        mu.Unlock()

        // Generate the new link in the format "new link name + .byteLink.bitly"
        newLink := fmt.Sprintf("%s.byteLink.bitly", customName)

        // Prepare the response
        response := fmt.Sprintf("Your new link is: <a href='%s' target='_blank'>%s</a>", newLink, newLink)
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, response)
        return
    }

    // If not a POST request, show an error
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// redirectLink handles redirection based on the custom name
func redirectLink(w http.ResponseWriter, r *http.Request) {
    customName := r.URL.Path[len("/bitly/"):] // Extract custom name from URL

    mu.Lock()
    link, exists := linkMap[customName]
    mu.Unlock()

    if exists {
        http.Redirect(w, r, link, http.StatusFound) // Redirect to the actual link
    } else {
        http.Error(w, "Link not found", http.StatusNotFound)
    }
}

func main() {
    // Serve the HTML page at the root URL
    http.HandleFunc("/", serveHTML)

    // Handle the form submission
    http.HandleFunc("/generate-free-link", generateLink)

    // Handle redirection for generated links
    http.HandleFunc("/bitly/", redirectLink)

    // Get the port from the environment variable or default to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port for local development
    }

    // Start the server
    fmt.Printf("Server is running on port %s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
