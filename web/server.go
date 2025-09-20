package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	// Custom handler to serve index.html for root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
		// For other paths, serve static files
		http.FileServer(http.Dir("")).ServeHTTP(w, r)
	})

	port := "3000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("🚀 Jira MCP Test UI Server starting on http://localhost:%s\n", port)
	fmt.Printf("📁 Serving files from: %s\n", "")
	fmt.Println("💡 Make sure your Jira MCP server is running on http://localhost:8080")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
