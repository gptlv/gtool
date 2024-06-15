package ad

import ldap "github.com/go-ldap/ldap/v3"

type AdService interface {
	GetByEmail(email string) (*ldap.Entry, error)
	GetUserGroups(user *ldap.Entry) []string
	ExtractCNFromDN(dn string) (string, error)
	RemoveUserFromGroup(user, group *ldap.Entry) (*ldap.Entry, error)
}
