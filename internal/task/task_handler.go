package task

import (
	"errors"
	"fmt"
	"log"
	"main/internal/dismissal"
	"main/internal/issue"
	"main/internal/object"

	"github.com/andygrunwald/go-jira"
)

type TaskHandler struct {
	issueService     issue.IssueService
	objectService    object.ObjectService
	dismissalService dismissal.DismissalService
}

func NewTaskHandler(is *issue.IssueService, os *object.ObjectService, ds *dismissal.DismissalService) *TaskHandler {
	return &TaskHandler{issueService: *is, objectService: *os, dismissalService: *ds}
}

func (h *TaskHandler) PrintLaptopDescription() error {
	var email string

	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	user, err := h.objectService.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	res, err := h.objectService.GetUserLaptops(user)
	if err != nil {
		log.Fatal(err)
	}

	for _, laptop := range res.ObjectEntries {
		description, err := h.objectService.GetLaptopDescription(&laptop)
		if err != nil {
			return fmt.Errorf("failed to get %v laptop description: %w", laptop.ObjectKey, err)
		}

		fmt.Printf("\n Ноутбук %s\n %s \n %s \n\n %s \n\n", description.Name, description.ISC, description.Serial, description.Cost)
	}

	return nil
}

func (h *TaskHandler) DeactivateInsight() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Деактивировать в Insight"`

	deactivationIssues, err := h.issueService.GetAll(jql)
	if err != nil {
		return err
	}

	if len(deactivationIssues) == 0 {
		return errors.New("no deactivation issues")
	}

	var unresolved []jira.Issue

	for _, di := range deactivationIssues {
		fmt.Print("found an issue: ")
		h.issueService.PrintIssue(&di)
		parentIssue, err := h.issueService.GetByID(di.Fields.Parent.ID)
		if err != nil {
			return err
		}

		fmt.Print("found a parent issue: ")
		h.issueService.PrintIssue(parentIssue)

		ss, err := h.issueService.GetSubtaskByComponent(parentIssue, "Возврат оборудования")
		if err != nil {
			panic(err)
		}

		if ss.Fields.Status.Name != "Closed" {
			fmt.Printf("parent issue %v has an incomplete subshipment task\n", parentIssue.Key)
			//block deactivation issue by the aforementioned subtask

			ds, err := h.issueService.GetSubtaskByComponent(parentIssue, "Insight")
			if err != nil {
				return err
			}

			di, err := h.issueService.GetByID(ds.ID)
			if err != nil {
				return err
			}

			bi, err := h.issueService.GetByID(ss.ID)
			if err != nil {
				return err
			}

			fmt.Printf("blocking %v by %v\n", di.Key, bi.Key)
			_, err = h.issueService.BlockByIssue(di, bi)
			if err != nil {
				return err
			}

			continue
		}

		issue, err := h.issueService.GetByID(parentIssue.ID)
		if err != nil {
			return fmt.Errorf("failed to get issue by ID: %w", err)
		}

		userEmail, err := h.issueService.GetUserEmail(issue)
		if err != nil {
			return err
		}

		fmt.Printf("found a user %v\n", userEmail)

		user, err := h.objectService.GetUserByEmail(userEmail)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		if user == nil {
			fmt.Printf("couldn't find user %v", userEmail)
			unresolved = append(unresolved, di)
			continue
		}

		getUserLaptopsRes, err := h.objectService.GetUserLaptops(user)
		if err != nil {
			return err
		}

		laptops := getUserLaptopsRes.ObjectEntries

		fmt.Printf("user %v has %v laptops\n", userEmail, len(laptops))

		var category string

		if len(laptops) > 0 {
			category = "Corporate laptop"
		} else {
			category = "BYOD"
		}

		fmt.Printf("changing %v's status to %v\n", user.ObjectKey, category)
		_, err = h.objectService.SetUserCategory(user, category)
		if err != nil {
			return fmt.Errorf("failed to set user category: %w", err)
		}

		fmt.Printf("disabling %v\n", user.ObjectKey)
		_, err = h.objectService.DisableUser(user)
		if err != nil {
			return fmt.Errorf("failed to disable user: %w", err)
		}

		deactivationSubtask, err := h.issueService.GetSubtaskByComponent(parentIssue, "Insight")
		if err != nil {
			return fmt.Errorf("failed to get subtask by component: %w", err)
		}

		deactivationIssue, err := h.issueService.GetByID(deactivationSubtask.ID)
		if err != nil {
			return fmt.Errorf("failed to get deactivation issue by ID: %w", err)
		}

		_, err = h.issueService.Close(deactivationIssue)
		if err != nil {
			return fmt.Errorf("failed to close deactivation issue: %w", err)
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

func (h *TaskHandler) GenerateDismissalRecords() error {
	csv, err := h.dismissalService.ReadCsvFile("info.csv")
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	var dismissalRecords []*dismissal.DismissalRecord

	for i, row := range csv {
		if i == 0 {
			continue
		}

		record, err := h.dismissalService.CreateDismissalRecord(row)
		if err != nil {
			return fmt.Errorf("failed to create dismissal record: %w", err)
		}

		isc := row[1]

		laptop, err := h.objectService.GetByISC(isc)
		if err != nil {
			return fmt.Errorf("failed to get laptop by isc: %w", err)
		}

		description, err := h.objectService.GetLaptopDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		record.Name = description.Name
		record.InventoryID = description.InventoryID
		record.Serial = description.Serial

		dismissalRecords = append(dismissalRecords, record)

	}

	for _, record := range dismissalRecords {
		templateNames := []string{"commitee", "dismissal", "record"}

		dirPath, err := h.dismissalService.CreateOutputDirectory(record.ISC)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		for _, templateName := range templateNames {
			template, err := h.dismissalService.CreateTemplate(record, templateName)
			if err != nil {
				return fmt.Errorf("failed to create a template %v: %w", templateName, err)
			}

			object, err := h.dismissalService.CreateObjectFromTemplate(template)
			if err != nil {
				return fmt.Errorf("failed to create object from template %v: %w", templateName, err)
			}

			file, err := h.dismissalService.CreateFile(dirPath, templateName, "pdf")
			if err != nil {
				return fmt.Errorf("failed to create file")
			}
			defer file.Close()

			err = h.dismissalService.CreatePDF(object, file)
			if err != nil {
				return fmt.Errorf("failed to generate pdf: %w", err)
			}
		}
	}

	return nil
}
