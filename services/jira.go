package services

import (
	"log"
	"sync"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

var JiraClient = sync.OnceValue[*jira.Client](func() *jira.Client {
	host, pat := loadJiraCredentials()

	if host == "" || pat == "" {
		log.Fatal("JIRA_HOST and JIRA_PAT are required")
	}

	// Use Bearer authentication for PAT tokens (self-hosted Jira)
	tp := jira.BearerAuthTransport{
		Token: pat,
	}

	instance, err := jira.NewClient(tp.Client(), host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create jira client"))
	}

	return instance
})
