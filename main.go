package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hdbrzgr/jira-mcp/v2/tools"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	envFile := flag.String("env", "", "Path to environment file (optional when environment variables are set directly)")
	httpPort := flag.String("http_port", "", "Port for HTTP server. If not provided, will use stdio")
	flag.Parse()

	// Load environment file if specified
	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			fmt.Printf("⚠️  Warning: Error loading env file %s: %v\n", *envFile, err)
		} else {
			fmt.Printf("✅ Loaded environment variables from %s\n", *envFile)
		}
	}

	// Check required environment variables
	host := os.Getenv("JIRA_HOST")
	pat := os.Getenv("JIRA_PAT")
	username := os.Getenv("JIRA_USERNAME")
	password := os.Getenv("JIRA_PASSWORD")

	missingEnvs := []string{}

	// JIRA_HOST is always required
	if host == "" {
		missingEnvs = append(missingEnvs, "JIRA_HOST")
	}

	// Check authentication: either PAT or username/password
	hasPAT := pat != ""
	hasBasicAuth := username != "" && password != ""

	if !hasPAT && !hasBasicAuth {
		if username == "" {
			missingEnvs = append(missingEnvs, "JIRA_USERNAME (for basic auth)")
		}
		if password == "" {
			missingEnvs = append(missingEnvs, "JIRA_PASSWORD (for basic auth)")
		}
		if pat == "" {
			missingEnvs = append(missingEnvs, "JIRA_PAT (for token auth)")
		}
	}

	if len(missingEnvs) > 0 {
		fmt.Println("❌ Configuration Error: Missing required environment variables")
		fmt.Println()
		fmt.Println("Missing variables:")
		for _, env := range missingEnvs {
			fmt.Printf("  - %s\n", env)
		}
		fmt.Println()
		fmt.Println("📋 Setup Instructions:")
		fmt.Println("Choose one of the following authentication methods:")
		fmt.Println()
		fmt.Println("🔑 Method 1: Personal Access Token (PAT) - For newer Jira versions")
		fmt.Println("1. Get your Personal Access Token (PAT) from your Jira instance")
		fmt.Println("   - Go to Jira > Settings > Personal Access Tokens")
		fmt.Println("   - Create a new token with appropriate permissions")
		fmt.Println("2. Set the environment variables:")
		fmt.Println("   JIRA_HOST=http://localhost:8080")
		fmt.Println("   JIRA_PAT=your-personal-access-token")
		fmt.Println()
		fmt.Println("🔑 Method 2: Username/Password - For older Jira versions (v2 API)")
		fmt.Println("1. Use your Jira username and password")
		fmt.Println("2. Set the environment variables:")
		fmt.Println("   JIRA_HOST=http://localhost:8080")
		fmt.Println("   JIRA_USERNAME=your-username")
		fmt.Println("   JIRA_PASSWORD=your-password")
		fmt.Println()
		fmt.Println("📁 Configuration Options:")
		fmt.Println("   Option A - Using .env file:")
		fmt.Println("   Create a .env file with one of the authentication methods above")
		fmt.Println()
		fmt.Println("   Option B - Using environment variables:")
		fmt.Println("   export JIRA_HOST=http://localhost:8080")
		fmt.Println("   # For PAT:")
		fmt.Println("   export JIRA_PAT=your-personal-access-token")
		fmt.Println("   # OR for username/password:")
		fmt.Println("   export JIRA_USERNAME=your-username")
		fmt.Println("   export JIRA_PASSWORD=your-password")
		fmt.Println()
		fmt.Println("   Option C - Using Docker:")
		fmt.Println("   # For PAT:")
		fmt.Printf("   docker run -e JIRA_HOST=http://localhost:8080 \\\n")
		fmt.Printf("              -e JIRA_PAT=your-personal-access-token \\\n")
		fmt.Printf("              ghcr.io/nguyenvanduocit/jira-mcp:latest\n")
		fmt.Println("   # For username/password:")
		fmt.Printf("   docker run -e JIRA_HOST=http://localhost:8080 \\\n")
		fmt.Printf("              -e JIRA_USERNAME=your-username \\\n")
		fmt.Printf("              -e JIRA_PASSWORD=your-password \\\n")
		fmt.Printf("              ghcr.io/nguyenvanduocit/jira-mcp:latest\n")
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("✅ All required environment variables are set")
	fmt.Printf("🔗 Connected to: %s\n", host)

	// Show which authentication method is being used
	if hasPAT {
		fmt.Println("🔑 Using Personal Access Token (PAT) authentication")
	} else {
		fmt.Println("🔑 Using Username/Password (Basic) authentication")
	}

	mcpServer := server.NewMCPServer(
		"Jira MCP",
		"1.0.1",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithRecovery(),
	)

	// Register available Jira tools
	tools.RegisterJiraIssueTool(mcpServer)
	tools.RegisterJiraSearchTool(mcpServer)
	tools.RegisterJiraTransitionTool(mcpServer)
	tools.RegisterJiraCommentTools(mcpServer)
	tools.RegisterJiraSprintTool(mcpServer) // Disabled - returns empty registration
	// Temporarily disabled during migration to andygrunwald/go-jira:
	// tools.RegisterJiraStatusTool(mcpServer)
	// tools.RegisterJiraWorklogTool(mcpServer)
	// tools.RegisterJiraHistoryTool(mcpServer)
	// tools.RegisterJiraRelationshipTool(mcpServer)

	if *httpPort != "" {
		fmt.Println()
		fmt.Println("🚀 Starting Jira MCP Server in HTTP mode...")
		fmt.Printf("📡 Server will be available at: http://localhost:%s/mcp\n", *httpPort)
		fmt.Println()
		fmt.Println("📋 Cursor Configuration:")
		fmt.Println("Add the following to your Cursor MCP settings (.cursor/mcp.json):")
		fmt.Println()
		fmt.Println("```json")
		fmt.Println("{")
		fmt.Println("  \"mcpServers\": {")
		fmt.Println("    \"jira\": {")
		fmt.Printf("      \"url\": \"http://localhost:%s/mcp\"\n", *httpPort)
		fmt.Println("    }")
		fmt.Println("  }")
		fmt.Println("}")
		fmt.Println("```")
		fmt.Println()
		fmt.Println("💡 Tips:")
		fmt.Println("- Restart Cursor after adding the configuration")
		fmt.Println("- Test the connection by asking Claude: 'List my Jira projects'")
		fmt.Println("- Use '@jira' in Cursor to reference Jira-related context")
		fmt.Println()
		fmt.Println("🔄 Server starting...")

		// Create MCP server with CORS support and stateless mode for easier testing
		httpServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithEndpointPath("/mcp"),
			server.WithStateLess(true))

		// Create custom HTTP server with CORS middleware
		customServer := &http.Server{
			Addr:    fmt.Sprintf(":%s", *httpPort),
			Handler: corsMiddleware(httpServer),
		}

		if err := customServer.ListenAndServe(); err != nil && !isContextCanceled(err) {
			log.Fatalf("❌ Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil && !isContextCanceled(err) {
			log.Fatalf("❌ Server error: %v", err)
		}
	}
}

// corsMiddleware adds CORS headers to allow cross-origin requests from the test UI
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Session-ID")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// IsContextCanceled checks if the error is related to context cancellation
func isContextCanceled(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's directly context.Canceled
	if errors.Is(err, context.Canceled) {
		return true
	}

	// Check if the error message contains context canceled
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "context canceled") ||
		strings.Contains(errMsg, "operation was canceled") ||
		strings.Contains(errMsg, "context deadline exceeded")
}
