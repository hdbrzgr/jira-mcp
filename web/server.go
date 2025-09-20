package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get the directory where this server.go file is located
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	webDir := filepath.Dir(execPath)

	// If running with go run, use the current directory
	if filepath.Base(webDir) == "go-build" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		webDir = wd
	}

	// Custom handler to serve index.html for root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}
		// For other paths, serve static files
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
	})

	port := "3000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("🚀 Jira MCP Test UI Server starting on http://localhost:%s\n", port)
	fmt.Printf("📁 Serving files from: %s\n", webDir)
	fmt.Println("💡 Make sure your Jira MCP server is running on http://localhost:8080")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
