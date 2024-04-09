package tasks

import (
	"errors"
	"fmt"
	"main/insight"
	"main/issue"
	"os"

	"github.com/andygrunwald/go-jira"
)

type LaptopDescription struct {
	Name   string
	ISC    string
	Serial string
}

func GetLaptopDescription(client *jira.Client, email string) {
	laptops, err := insight.GetUserLaptops(client, email)
	if err != nil {
		panic(err)
	}

	if len(laptops.ObjectEntries) > 1 {
		panic(errors.New("user has more than one laptop"))
	}

	serial := laptops.ObjectEntries[0].Label
	isc := laptops.ObjectEntries[0].ObjectKey
	name := laptops.ObjectEntries[0].Attributes[1].ObjectAttributeValues[0].Value

	description := LaptopDescription{
		name, isc, serial,
	}

	fmt.Println(description)

	for _, attribute := range laptops.ObjectEntries[0].Attributes {
		if attribute.ObjectTypeAttributeID == 4184 {
			fmt.Println(attribute.ObjectAttributeValues[0].Value)
		}
	}

}

func DeactivateInsight(client *jira.Client) {
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

		fmt.Printf("user %v has %v laptops\n", userEmail, len(laptops.ObjectEntries))

		var category string

		if len(laptops.ObjectEntries) > 0 {
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
