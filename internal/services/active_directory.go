package services

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
	"github.com/gptlv/gtools/internal/interfaces"
)

type activeDirectoryService struct {
	conn *ldap.Conn
}

func NewActiveDirectoryService(conn *ldap.Conn) interfaces.ActiveDirectoryService {
	return &activeDirectoryService{conn: conn}
}

func (s *activeDirectoryService) GetByCN(cn string) (*ldap.Entry, error) {
	if cn == "" {
		return nil, errors.New("empty cn")
	}

	baseDN := os.Getenv("LDAP_BASE_DN")
	filter := fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(cn))

	searchGroupReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName", "cn"}, []ldap.Control{})

	groupRes, err := s.conn.Search(searchGroupReq)
	if err != nil {
		log.Fatal("failed to query LDAP: %w", err)
	}

	if len(groupRes.Entries) == 0 {
		return nil, errors.New("no such entry")
	}

	group := groupRes.Entries[0]

	return group, nil

}

func (s *activeDirectoryService) GetByEmail(email string) (*ldap.Entry, error) {
	if email == "" {
		return nil, errors.New("empty email")
	}

	baseDN := os.Getenv("LDAP_BASE_DN")
	filter := fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(email))

	searchUserReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName", "memberOf", "mail", "userAccountControl"}, []ldap.Control{})

	userRes, err := s.conn.Search(searchUserReq)
	if err != nil {
		return nil, fmt.Errorf("failed to query LDAP: %w", err)
	}
	if len(userRes.Entries) > 1 {
		return nil, fmt.Errorf("found multiple accounts for %v", email)
	}

	user := userRes.Entries[0]

	return user, nil
}

func (s *activeDirectoryService) GetUserGroups(user *ldap.Entry) []string {
	return user.GetAttributeValues("memberOf")
}

func (s *activeDirectoryService) ExtractCNFromDN(dn string) (string, error) {
	parsedDN, err := ldap.ParseDN(dn)
	if err != nil {
		return "", err
	}

	for _, rdn := range parsedDN.RDNs {
		for _, atv := range rdn.Attributes {
			if atv.Type == "CN" || atv.Type == "cn" {
				return atv.Value, nil
			}
		}
	}
	return "", fmt.Errorf("CN not found in DN: %s", dn)
}

func (s *activeDirectoryService) RemoveUserFromGroup(user, group *ldap.Entry) (*ldap.Entry, error) {
	modify := ldap.NewModifyRequest(group.DN, []ldap.Control{})
	modify.Delete("member", []string{user.DN})

	err := s.conn.Modify(modify)
	if err != nil {
		return nil, err
	}

	email := user.GetAttributeValue("mail")

	updatedUser, err := s.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return updatedUser, nil
}

func (s *activeDirectoryService) AddUserToGroup(user, group *ldap.Entry) (*ldap.Entry, error) {
	if user == nil {
		return nil, errors.New("empty user")

	}

	if group == nil {
		return nil, errors.New("empty group")
	}

	modify := ldap.NewModifyRequest(group.DN, []ldap.Control{})
	modify.Add("member", []string{user.DN})

	err := s.conn.Modify(modify)
	if err != nil {
		return nil, err
	}

	email := user.GetAttributeValue("mail")

	updatedUser, err := s.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return updatedUser, nil

}

func (s *activeDirectoryService) UpdateDN(user *ldap.Entry, newSup string) (*ldap.Entry, error) {
	cn := user.GetAttributeValue("cn")
	rdn := fmt.Sprintf("CN=%v", cn)

	req := ldap.NewModifyDNRequest(user.DN, rdn, true, newSup)
	err := s.conn.ModifyDN(req)
	if err != nil {
		return nil, fmt.Errorf("failed to modify userDN: %s", err)
	}

	updatedUser, err := s.GetByCN(cn)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil

}
