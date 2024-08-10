package common

import (
	"main/config"

	"github.com/andygrunwald/go-jira"
)

func GetJiraClient() (*jira.Client, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	tp := jira.BearerAuthTransport{
		Token: cfg.Jira.Token,
	}

	client, err := jira.NewClient(tp.Client(), cfg.Jira.URL)
	if err != nil {
		return nil, err
	}

	return client, nil
}
