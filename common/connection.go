package common

import (
	"main/config"

	"github.com/go-ldap/ldap/v3"
)

func GetLDAPConnection() (*ldap.Conn, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	conn, err := ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.Bind(cfg.ActiveDirectory.AdminDN, cfg.ActiveDirectory.AdminPassword)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
