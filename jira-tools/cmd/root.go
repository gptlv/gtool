package cmd

import (
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/go-ldap/ldap/v3"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var Client *jira.Client
var Conn *ldap.Conn

func init() {
	if envErr := godotenv.Load(".env"); envErr != nil {
		fmt.Println(".env file missing")
	}
}

var RootCmd = &cobra.Command{
	Use:   "jira-tools",
	Short: "CLI application for mundane tasks in jira",
	Long:  "A CLI application built with Cobra for various operations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func Execute() {
	injectConnections()

	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
