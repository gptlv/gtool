package cmd

import (
	"log"
	"main/config"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/go-ldap/ldap/v3"
)

func injectConnections() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	tp := jira.BearerAuthTransport{
		Token: cfg.Jira.Token,
	}

	Client, err = jira.NewClient(tp.Client(), cfg.Jira.URL)
	if err != nil {
		log.Fatal(err)
	}

	Conn, err = ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer Conn.Close()

	err = Conn.Bind(os.Getenv("ADMIN_DN"), os.Getenv("ADMIN_PASS"))
	if err != nil {
		log.Fatal(err)
	}

}
