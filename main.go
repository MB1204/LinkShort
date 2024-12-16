package main

import (
    "fmt"
    "net/http"
    "os"
)

// serveHTML serves the HTML content from the index.html file
func serveHTML(w http.ResponseWriter, r *http.Request) {
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

        // Generate the new link
        newLink := fmt.Sprintf("https://bytelink.com/bitly/%s", customName)

        // Prepare the response
        response := fmt.Sprintf("Your new link is: <a href='%s' target='_blank'>%s</a>", newLink, newLink)
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, response)
        return
    }

    // If not a POST request, show an error
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func main() {
    // Serve the HTML page at the root URL
    http.HandleFunc("/", serveHTML)

    // Handle the form submission
    http.HandleFunc("/generate-link", generateLink)

    // Start the server
    fmt.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }
}