package main

import (
	"fmt"
	"log"
	"main/internal/dismissal"
	"main/internal/issue"
	"main/internal/object"
	"main/internal/task"
	"os"
	"strconv"

	pdf "github.com/adrg/go-wkhtmltopdf"
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

	is := issue.NewIssueService(client)
	os := object.NewObjectService(client)
	ds := dismissal.NewDismissalService()

	th := task.NewTaskHandler(&is, &os, &ds)

	fmt.Print("\033[H\033[2J")
	fmt.Println("1) Deactivate insight")
	fmt.Println("2) Generate dismissal documents")
	fmt.Println("3) Get laptop description")
	fmt.Println("4) Assign all deactivate insight issues to me")
	fmt.Println("5) Show issues with empty component")
	fmt.Println("6) Update block trainee cc issue")

	var n int
	for {
		var input string
		fmt.Print("Your choice: ")
		fmt.Scanln(&input)

		n, err = strconv.Atoi(input)
		if err == nil && (1 <= n && n <= 6) {
			break
		}
		fmt.Println("Invalid choice.")
	}

	if n == 1 {
		err := th.DeactivateInsight()
		if err != nil {
			panic(err)
		}
	}

	if n == 2 {
		if err := pdf.Init(); err != nil {
			log.Fatal(fmt.Errorf("failed to initialize pdf: %w", err))
		}
		defer pdf.Destroy()

		err := th.GenerateDismissalRecords()
		if err != nil {
			panic(err)
		}
	}

	if n == 3 {
		err := th.PrintLaptopDescription()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 4 {
		err := th.AssignAllDeactivateInsightIssuesToMe()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 5 {
		err := th.ShowIssuesWithEmptyComponent()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 6 {
		err := th.UpdateBlockTraineeIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

}
