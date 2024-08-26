package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
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

	arg.MustParse(&args)

	if args.Issue != nil {

		if args.Issue.Action == "process-insight" {
			if err := gt.ProcessInsight(); err != nil {
				log.Fatal(err)
			}
		}

		if args.Issue.Action == "process-ldap" {
			if err := gt.ProcessLDAP(); err != nil {
				log.Fatal(err)
			}
		}

		if args.Issue.Action == "process-staff" {
			if err := gt.ProcessStaff(args.Issue.Component); err != nil {
				log.Fatal(err)
			}
		}

		if args.Issue.Action == "grant-access" {
			if err := gt.GrantAccess(args.Issue.Key); err != nil {
				log.Fatal(err)
			}
		}

		if args.Issue.Action == "assign" {
			if err := gt.AssignAll(args.Issue.Component); err != nil {
				log.Fatal(err)
			}

		}

		if args.Issue.Action == "update-trainee" {
			if err := gt.UpdateBlockTraineeIssue(args.Issue.Key); err != nil {
				log.Fatal(err)
			}
		}

		if args.Issue.Action == "show-empty" {
			if err := gt.ShowEmpty(); err != nil {
				log.Fatal(err)
			}
		}

	}

	if args.Asset != nil {

		if args.Asset.Action == "print-description" {
			if err := gt.PrintDescription(args.Asset.Email); err != nil {
				log.Fatal(err)
			}
		}

		if args.Asset.Action == "generate-records" {
			if err := gt.GenerateRecords(args.Asset.StartID); err != nil {
				log.Fatal(err)
			}
		}
	}

	if args.LDAP != nil {

		if args.LDAP.Action == "add-group" {
			if err := gt.AddGroup(args.LDAP.Emails, args.LDAP.CNs); err != nil {
				log.Fatal(err)
			}
		}

	}
}
