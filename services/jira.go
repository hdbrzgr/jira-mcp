package services

import (
	"log"
	"sync"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

var JiraClient = sync.OnceValue[*jira.Client](func() *jira.Client {
	config := loadJiraCredentials()

	var instance *jira.Client
	var err error

	if config.UseBasicAuth {
		// Use Basic authentication for username/password (older Jira versions)
		log.Println("Using Basic authentication (username/password)")
		tp := jira.BasicAuthTransport{
			Username: config.Username,
			Password: config.Password,
		}
		instance, err = jira.NewClient(tp.Client(), config.Host)
	} else {
		// Use Bearer authentication for PAT tokens (newer Jira versions)
		log.Println("Using Bearer authentication (PAT)")
		tp := jira.BearerAuthTransport{
			Token: config.PAT,
		}
		instance, err = jira.NewClient(tp.Client(), config.Host)
	}

	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create jira client"))
	}

	return instance
})
