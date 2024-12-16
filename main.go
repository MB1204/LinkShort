package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

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

        // Prepare the JSON response
        response := map[string]string{"newLink": newLink}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    // If not a POST request, show an error
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func main() {
    // Handle the form submission
    http.HandleFunc("/generate-link", generateLink)

    // Start the server
    fmt.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}