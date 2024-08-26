package main

type IssueCmd struct {
	Action    string `arg:"positional,required"`
	Key       string
	Component string
}

type AssetCmd struct {
	Action  string `arg:"positional,required"`
	Email   string
	StartID int `arg:"--start-id" default:"1"`
}

type LDAPCmd struct {
	Action string `arg:"positional,required"`
	Emails []string
	CNs    []string
}

var args struct {
	Issue *IssueCmd `arg:"subcommand:issue"`
	Asset *AssetCmd `arg:"subcommand:asset"`
	LDAP  *LDAPCmd  `arg:"subcommand:ldap"`
}
