package tasks

import (
	"errors"
	"fmt"
	"main/dismissal"
	"main/insight"
	"main/issue"

	"github.com/andygrunwald/go-jira"
)

type LaptopDescription struct {
	Name   string
	Isc    string
	Serial string
	Cost   string
}

func GetUserLaptopDescription(client *jira.Client) error {
	const ISC_ATTRIBUTE_ID = 879
	const NAME_ATTRIBUTE_ID = 880
	const SERIAL_ATTRIBUTE_ID = 889
	const COST_ATTRIBUTE_ID = 4184

	var email string
	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	if email == "" {
		panic(errors.New("empty email"))
	}

	user, err := insight.GetUserByEmail(client, email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	laptops, err := insight.GetUserLaptops(client, user)
	if err != nil {
		return fmt.Errorf("failed to get user laptops: %w", err)
	}

	if len(laptops.ObjectEntries) > 1 {
		fmt.Printf("!!!\nuser has more than one laptop\n!!!\n")
	}

	var name, isc, serial, cost string
	var res []LaptopDescription

	for _, entry := range laptops.ObjectEntries {
		d := new(LaptopDescription)
		for _, attribute := range entry.Attributes {
			attributeValue := attribute.ObjectAttributeValues[0].Value

			switch attribute.ObjectTypeAttributeID {
			case NAME_ATTRIBUTE_ID:
				name = attributeValue
				d.Name = name
			case ISC_ATTRIBUTE_ID:
				isc = attributeValue
				d.Isc = isc
			case SERIAL_ATTRIBUTE_ID:
				serial = attributeValue
				d.Serial = serial
			case COST_ATTRIBUTE_ID:
				cost = attributeValue
				d.Cost = cost
			}
		}
		res = append(res, *d)
	}

	PrintLaptopDescription(res)

	return nil
	// return res, nil

}

func PrintLaptopDescription(description []LaptopDescription) {
	for _, d := range description {
		fmt.Printf("\nНоутбук %s\n %s \n %s \n\n %s \n\n", d.Name, d.Isc, d.Serial, d.Cost)
	}
}

func DeactivateInsight(client *jira.Client) error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Деактивировать в Insight"`

	deactivationIssues, err := issue.GetAll(client, jql)
	if err != nil {
		return err
	}

	if len(deactivationIssues) == 0 {
		return errors.New("no deactivation issues")
	}

	var parentIssues []jira.Issue

	for _, di := range deactivationIssues {
		fmt.Print("found an issue: ")
		issue.PrintIssue(&di)
		parentIssue, err := issue.GetByID(client, di.Fields.Parent.ID)
		if err != nil {
			return err
		}

		fmt.Print("found a parent issue: ")
		issue.PrintIssue(parentIssue)
		parentIssues = append(parentIssues, *parentIssue)
	}

	for _, pi := range parentIssues {
		ss, err := issue.GetSubtaskByComponent(client, &pi, "Возврат оборудования")
		if err != nil {
			panic(err)
		}

		if ss.Fields.Status.Name != "Closed" {
			fmt.Printf("parent issue %v has an incomplete subshipment task\n", pi.Key)
			//block deactivation issue by the aforementioned subtask

			ds, err := issue.GetSubtaskByComponent(client, &pi, "Insight")
			if err != nil {
				return err
			}

			di, err := issue.GetByID(client, ds.ID)
			if err != nil {
				return err
			}

			bi, err := issue.GetByID(client, ss.ID)
			if err != nil {
				return err
			}

			fmt.Printf("blocking %v by %v\n", di.Key, bi.Key)
			_, err = issue.BlockByIssue(client, di, bi)
			if err != nil {
				return err
			}

			continue
		}

		userEmail, err := issue.GetUserEmail(client, pi.Key)
		if err != nil {
			return err
		}

		fmt.Printf("found a user %v\n", userEmail)

		user, err := insight.GetUserByEmail(client, userEmail)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		laptops, err := insight.GetUserLaptops(client, user)
		if err != nil {
			return err
		}

		fmt.Printf("user %v has %v laptops\n", userEmail, len(laptops.ObjectEntries))

		var category string

		if len(laptops.ObjectEntries) > 0 {
			category = "Corporate laptop"
		} else {
			category = "BYOD"
		}

		// user, err := insight.GetUserByEmail(client, userEmail)
		// if err != nil {
		// 	panic(err)
		// }

		fmt.Printf("changing %v's status to %v\n", user.ObjectKey, category)
		_, err = insight.SetUserCategory(client, user, category)
		if err != nil {
			panic(err)
		}

		fmt.Printf("disabling %v\n", user.ObjectKey)
		_, err = insight.DisableUser(client, user)
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

	return nil
}

func GenerateDismissalDocuments() {
	const boss = "Козлов А.Ю."
	const lead = "Степанов А.С."
	doc := &dismissal.Document{
		Template: "record.html",
		Serial:   "1337",
		Name:     "макбук))",
		Isc:      "isc-228",
		Date:     "14 Апреля 2024",
		ID:       250,
		Boss:     boss,
		Lead:     lead,
		Flaws:    "none",
		Decision: "burn",
	}

	// commonDocument := &dismissal.Document{
	// 	Serial: "1337",
	// 	Name:   "laptop123",
	// 	Isc:    "isc-228",
	// 	Date:   "20.03.2000",
	// 	Boss:   "Козлов А.Ю.",
	// 	Lead:   "Степанов А.С.",
	// }

	template, err := dismissal.GenerateTemplate(doc)
	if err != nil {
		panic(err)
	}

	dismissal.GenerateDocument(template)
}
