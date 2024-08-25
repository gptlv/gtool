package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap"
	"github.com/joho/godotenv"
)

var config *Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load .env: %w", err))
	}

	config, err = NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}
}

func main() {
	tp := jira.BearerAuthTransport{
		Token: config.Jira.Token,
	}

	client, err := jira.NewClient(tp.Client(), config.Jira.URL)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ldap.DialURL(config.LDAP.URL)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Bind(config.LDAP.AdminDN, config.LDAP.AdminPassword)
	if err != nil {
		log.Fatal(err)
	}

	gt := New(client, conn)

	err = gt.GenerateRecords()
	if err != nil {
		log.Fatal(err)
	}
}
