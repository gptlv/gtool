package interfaces

import "main/internal/entities"

type ActiveDirectoryService interface {
	GetByEmail(email string) (*entities.User, error)
	GetByCN(cn string) (*entities.User, error)
	GetUserGroups(user *entities.User) []string
	// RemoveUserFromGroup(user, group *ldap.Entry) (*ldap.Entry, error)
	// AddUserToGroup(user, group *ldap.Entry) (*ldap.Entry, error)
	ExtractCNFromDN(dn string) (string, error)
	UpdateDN(user *entities.User, newDN string) (*entities.User, error)
}

// type ActiveDirectoryHandler interface {
// 	RemoveVPNGroupsFromUsers()
// 	MoveUsersToNewOU()
// }
