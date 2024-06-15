package ad

import (
	"errors"
	"fmt"
	"log"
	"os"

	ldap "github.com/go-ldap/ldap/v3"
)

type adService struct {
	conn *ldap.Conn
}

func NewAdService(conn *ldap.Conn) AdService {
	return &adService{conn: conn}
}

func (s *adService) GetByEmail(email string) (*ldap.Entry, error) {
	if email == "" {
		return nil, errors.New("empty email")
	}

	baseDN := os.Getenv("BASE_DN")
	filter := fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(email))

	searchUserReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName", "memberOf", "mail"}, []ldap.Control{})

	userRes, err := s.conn.Search(searchUserReq)
	if err != nil {
		log.Fatal("failed to query LDAP: %w", err)
	}

	user := userRes.Entries[0]

	return user, nil
}

func (s *adService) GetUserGroups(user *ldap.Entry) []string {
	return user.GetAttributeValues("memberOf")
}

func (s *adService) ExtractCNFromDN(dn string) (string, error) {
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

func (s *adService) RemoveUserFromGroup(user, group *ldap.Entry) (*ldap.Entry, error) {
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
