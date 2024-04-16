package main

import (
	"errors"
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

	var email string
	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	if email == "" {
		panic(errors.New("empty email"))
	}

	descriptions := tasks.GetLaptopDescription(client, email)
	tasks.PrintLaptopDescription(descriptions)

	// tasks.GenerateDismissalDocuments()

}
