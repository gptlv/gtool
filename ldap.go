package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap"
)

func (g *gtool) AddGroup(emails, groupCNs []string) error {
	var users []*ldap.Entry
	var groups []*ldap.Entry

	for _, email := range emails {
		log.Info(fmt.Sprintf("getting AD user by email: %s", email))
		query := fmt.Sprintf("mail=%s", email)
		user, err := g.searchEntry(query, []string{})
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		users = append(users, user)
	}

	for _, groupCN := range groupCNs {
		log.Info(fmt.Sprintf("getting AD group by cn: %s", groupCN))
		query := fmt.Sprintf("cn=%s", groupCN)
		group, err := g.searchEntry(query, []string{})
		if err != nil {
			return fmt.Errorf("failed to get group by cn: %w", err)
		}

		groups = append(groups, group)
	}

	for _, user := range users {
		for _, group := range groups {
			log.Info(fmt.Sprintf("adding user %s to group %s", user.GetAttributeValue("mail"), group.GetAttributeValue("cn")))
			err := g.addUserToGroup(user, group)
			if err != nil {
				log.Error(fmt.Errorf("failed to add user to group: %w", err))
			}
		}
	}

	return nil
}

func (g *gtool) searchEntry(query string, attributes []string) (*ldap.Entry, error) {
	searchEntryFilter := fmt.Sprintf("(%s)", ldap.EscapeFilter(query))
	searchEntryReq := ldap.NewSearchRequest(config.LDAP.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, searchEntryFilter, attributes, []ldap.Control{})
	searchUserResult, err := g.conn.Search(searchEntryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}

	if len(searchUserResult.Entries) < 1 {
		return nil, errors.New("no such entry")
	}

	if len(searchUserResult.Entries) > 1 {
		return nil, errors.New("found multiple entries")
	}

	return searchUserResult.Entries[0], nil
}

func (g *gtool) addUserToGroup(user, group *ldap.Entry) error {
	modify := ldap.NewModifyRequest(group.DN, []ldap.Control{})
	modify.Add("member", []string{user.DN})

	return g.conn.Modify(modify)
}
