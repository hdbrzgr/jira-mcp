package services

import (
	"log"
	"os"
)

// AuthConfig holds the authentication configuration for Jira
type AuthConfig struct {
	Host         string
	PAT          string
	Username     string
	Password     string
	UseBasicAuth bool
}

// loadJiraCredentials loads Jira credentials from environment variables
// Supports both PAT (Personal Access Token) and username/password authentication
// For older Jira versions (v2 API), use JIRA_USERNAME and JIRA_PASSWORD
// For newer versions, use JIRA_PAT
func loadJiraCredentials() AuthConfig {
	host := os.Getenv("JIRA_HOST")
	pat := os.Getenv("JIRA_PAT")
	username := os.Getenv("JIRA_USERNAME")
	password := os.Getenv("JIRA_PASSWORD")

	if host == "" {
		log.Fatal("JIRA_HOST is required, please set it in MCP Config")
	}

	// Check if we have PAT or username/password
	hasPAT := pat != ""
	hasBasicAuth := username != "" && password != ""

	if !hasPAT && !hasBasicAuth {
		log.Fatal("Either JIRA_PAT or both JIRA_USERNAME and JIRA_PASSWORD are required for authentication")
	}

	if hasPAT && hasBasicAuth {
		log.Println("Both PAT and username/password provided, using PAT authentication")
	}

	return AuthConfig{
		Host:         host,
		PAT:          pat,
		Username:     username,
		Password:     password,
		UseBasicAuth: !hasPAT && hasBasicAuth,
	}
}
