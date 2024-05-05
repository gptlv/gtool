package tasks

import (
	"errors"
	"fmt"
	"main/dismissal"
	"main/insight"
	"main/issue"
	"main/util"

	"github.com/andygrunwald/go-jira"
)

func GetUserLaptopDescription(client *jira.Client) error {
	var email string

	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	if email == "" {
		return errors.New("empty email")
	}

	user, err := insight.GetUserByEmail(client, email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	laptops, err := insight.GetUserLaptops(client, user)
	if err != nil {
		return fmt.Errorf("failed to get user laptops: %w", err)
	}

	if len(laptops) > 1 {
		fmt.Printf("!!!\nuser has more than one laptop\n!!!\n")
	}

	for _, laptop := range laptops {
		d, err := insight.GetLaptopDescription(client, &laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		err = insight.PrintLaptopDescription(d)
		if err != nil {
			return fmt.Errorf("failed to print laptop description: %w", err)
		}

	}

	return nil
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

	var unresolved []jira.Issue

	for _, di := range deactivationIssues {
		fmt.Print("found an issue: ")
		issue.PrintIssue(&di)
		parentIssue, err := issue.GetByID(client, di.Fields.Parent.ID)
		if err != nil {
			return err
		}

		fmt.Print("found a parent issue: ")
		issue.PrintIssue(parentIssue)

		ss, err := issue.GetSubtaskByComponent(client, parentIssue, "Возврат оборудования")
		if err != nil {
			panic(err)
		}

		if ss.Fields.Status.Name != "Closed" {
			fmt.Printf("parent issue %v has an incomplete subshipment task\n", parentIssue.Key)
			//block deactivation issue by the aforementioned subtask

			ds, err := issue.GetSubtaskByComponent(client, parentIssue, "Insight")
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

		userEmail, err := issue.GetUserEmail(client, parentIssue.Key)
		if err != nil {
			return err
		}

		fmt.Printf("found a user %v\n", userEmail)

		user, err := insight.GetUserByEmail(client, userEmail)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		if user == nil {
			fmt.Printf("couldn't find user %v", userEmail)
			unresolved = append(unresolved, di)
			continue
		}

		laptops, err := insight.GetUserLaptops(client, user)
		if err != nil {
			return err
		}

		fmt.Printf("user %v has %v laptops\n", userEmail, len(laptops))

		var category string

		if len(laptops) > 0 {
			category = "Corporate laptop"
		} else {
			category = "BYOD"
		}

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

		deactivationSubtask, err := issue.GetSubtaskByComponent(client, parentIssue, "Insight")
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

	if len(unresolved) > 0 {
		fmt.Print("unresolved issues:\n")
		for i, ui := range unresolved {
			fmt.Printf("%v. %v", i, ui.Key)
		}
	}

	return nil
}

func CreateDocuments(client *jira.Client, filePath string) error {
	csv, err := util.ReadCsvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	dismissalRecords, err := dismissal.CreateDismissalRecords(csv)
	if err != nil {
		return fmt.Errorf("failed to create dismissal records: %w", err)
	}

	for _, record := range dismissalRecords {
		templateNames := []string{"commitee", "dismissal", "record"}

		dismissalDocument, err := dismissal.CreateDismissalDocument(client, &record)
		if err != nil {
			return fmt.Errorf("failed to generate dismissal document: %w", err)
		}

		dirPath, err := util.CreateOutputDirectory(dismissalDocument.LaptopDescription.ISC)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		for _, templateName := range templateNames {
			template, err := dismissal.CreateTemplate(dismissalDocument, templateName)
			if err != nil {
				return fmt.Errorf("failed to create a template %v: %w", templateName, err)
			}

			object, err := dismissal.CreateObjectFromTemplate(template)
			if err != nil {
				return fmt.Errorf("failed to create object from template %v: %w", templateName, err)
			}

			file, err := util.CreateFile(dirPath, templateName, "pdf")
			if err != nil {
				return fmt.Errorf("failed to create file")
			}
			defer file.Close()

			err = dismissal.CreatePDF(object, file)
			if err != nil {
				return fmt.Errorf("failed to generate pdf: %w", err)
			}
		}
	}

	return nil
}
