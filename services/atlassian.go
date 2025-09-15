package services

import (
	"log"
	"os"
)

func loadJiraCredentials() (host, pat string) {
	host = os.Getenv("JIRA_HOST")
	pat = os.Getenv("JIRA_PAT")

	if host == "" || pat == "" {
		log.Fatal("JIRA_HOST and JIRA_PAT are required, please set them in MCP Config")
	}

	return host, pat
}
