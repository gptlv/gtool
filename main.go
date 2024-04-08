package main

import (
	"fmt"
	"main/insight"
	"main/issue"
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

	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Деактивировать в Insight"`

	deactivationIssues, err := issue.GetAll(client, jql)
	if err != nil {
		panic(err)
	}

	if len(deactivationIssues) == 0 {
		fmt.Println("no deactivation issues")
		os.Exit(1)
	}

	var parentIssues []jira.Issue

	for _, di := range deactivationIssues {
		parentIssue, err := issue.GetByID(client, di.Fields.Parent.ID)
		if err != nil {
			panic(err)
		}

		fmt.Printf("found a parent issue %v\n", parentIssue.Key)
		parentIssues = append(parentIssues, *parentIssue)
	}

	for _, pi := range parentIssues {
		ss, err := issue.GetSubtaskByComponent(client, &pi, "Возврат оборудования")
		if err != nil {
			panic(err)
		}

		if ss.Fields.Status.Name != "Closed" {
			fmt.Printf("parent issue %v has an incomplete subshipment task\n", pi.Key)
			fmt.Printf("continuing...\n")
			//block deactivation issue by the aforementioned subtask
			continue
		}

		userEmail, err := issue.GetUserEmail(client, pi.Key)
		if err != nil {
			panic(err)
		}

		fmt.Printf("found a user %v\n", userEmail)

		laptops, err := insight.GetUserLaptops(client, userEmail)
		if err != nil {
			panic(err)
		}

		fmt.Printf("user %v has %v laptops\n", userEmail, len(laptops))

		var category string

		if len(laptops) > 0 {
			category = "Corporate laptop"
		} else {
			category = "BYOD"
		}

		ISC, err := insight.GetUserISC(client, userEmail)
		if err != nil {
			panic(err)
		}

		fmt.Printf("changing %v's status to %v\n", ISC, category)
		_, err = insight.SetUserCategory(client, ISC, category)
		if err != nil {
			panic(err)
		}

		fmt.Printf("disabling %v\n", ISC)
		_, err = insight.DisableUser(client, ISC)
		if err != nil {
			panic(err)
		}

		deactivationSubtask, err := issue.GetSubtaskByComponent(client, &pi, "Insight")
		if err != nil {
			panic(err)
		}

		deactivationIssue, err := issue.GetByID(client, deactivationSubtask.ID)
		if err != nil {
			panic(err)
		}

		_, err = issue.Close(client, deactivationIssue)
		if err != nil {
			panic(err)
		}
	}
}
