package main

import (
	"fmt"
	"main/tasks"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
)

func init() {
	if envErr := godotenv.Load(".env"); envErr != nil {
		fmt.Println(".env file missing")
	}
}

func main() {
	tp := jira.BearerAuthTransport{
		Token: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		panic(err)
	}

	// tasks.DeactivateInsight(client)
	email := ""
	tasks.GetLaptopDescription(client, email)

}
