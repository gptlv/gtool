package main

import (
	"fmt"
	"log"
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
		log.Fatal(err)
	}

	err = tasks.DeactivateInsight(client)
	if err != nil {
		log.Fatal(err)
	}

	// err = tasks.GetUserLaptopDescription(client)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := pdf.Init(); err != nil {
	// 	log.Fatal(fmt.Errorf("failed to initialize pdf: %w", err))
	// }
	// defer pdf.Destroy()

	// err = tasks.CreateDocuments(client, "info.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
