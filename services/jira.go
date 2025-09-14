package services

import (
	"log"
	"sync"

	jira "github.com/ctreminiom/go-atlassian/jira/v2"
	"github.com/pkg/errors"
)

var JiraClient = sync.OnceValue[*jira.Client](func() *jira.Client {
	host, pat := loadJiraCredentials()

	if host == "" || pat == "" {
		log.Fatal("JIRA_HOST and JIRA_PAT are required")
	}

	instance, err := jira.New(nil, host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create jira client"))
	}

	// Use PAT authentication for local Jira 10.3.2
	instance.Auth.SetBearerToken(pat)

	return instance
})
