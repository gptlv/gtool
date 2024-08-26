package main

import (
	"github.com/andygrunwald/go-jira"
	"github.com/go-ldap/ldap"
)

type GTool interface {
	AssignAll(component string) error
	ProcessInsight() error
	ProcessLDAP() error
	ProcessStaff(component string) error
	GrantAccess(key string) error
	UpdateBlockTraineeIssue(key string) error
	ShowEmpty() error
	PrintDescription(email string) error
	GenerateRecords(startID int) error
	AddGroup(emails, groupCNs []string) error
}

type gtool struct {
	client *jira.Client
	conn   *ldap.Conn
}

func New(client *jira.Client, conn *ldap.Conn) GTool {
	return &gtool{client: client, conn: conn}
}
