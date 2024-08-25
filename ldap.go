package main

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap"
)

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
