package services

import (
	"log"
	"os"
	"sync"

	"github.com/ctreminiom/go-atlassian/jira/agile"
	"github.com/pkg/errors"
)

func loadJiraCredentials() (host, pat string) {
	host = os.Getenv("JIRA_HOST")
	pat = os.Getenv("JIRA_PAT")

	if host == "" || pat == "" {
		log.Fatal("JIRA_HOST and JIRA_PAT are required, please set them in MCP Config")
	}

	return host, pat
}

var AgileClient = sync.OnceValue[*agile.Client](func() *agile.Client {
	host, pat := loadJiraCredentials()

	instance, err := agile.New(nil, host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create agile client"))
	}

	// Use PAT authentication for local Jira 10.3.2
	instance.Auth.SetBearerToken(pat)

	return instance
})
