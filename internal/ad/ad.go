package ad

import ldap "github.com/go-ldap/ldap/v3"

type AdService interface {
	GetByEmail(email string) (*ldap.Entry, error)
	GetByCN(cn string) (*ldap.Entry, error)
	GetUserGroups(user *ldap.Entry) []string
	RemoveUserFromGroup(user, group *ldap.Entry) (*ldap.Entry, error)
	AddUserToGroup(user, group *ldap.Entry) (*ldap.Entry, error)
	ExtractCNFromDN(dn string) (string, error)
	UpdateDN(user *ldap.Entry, newDN string) (*ldap.Entry, error)
}
