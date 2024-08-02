package cmd

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/go-ldap/ldap/v3"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	if envErr := godotenv.Load(".env"); envErr != nil {
		fmt.Println(".env file missing")
	}
}

var rootCmd = &cobra.Command{
	Use:   "jira-tools",
	Short: "CLI application for mundane tasks in jira",
	Long:  "A CLI application built with Cobra for various operations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func Execute() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	tp := jira.BearerAuthTransport{
		Token: cfg.Jira.Token,
	}

	client, err := jira.NewClient(tp.Client(), cfg.Jira.URL)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// err = conn.Bind(os.Getenv("ADMIN_DN"), os.Getenv("ADMIN_PASS"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.WithValue(context.Background(), "client", client)
	ctx = context.WithValue(ctx, "conn", conn)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
