package task

import (
	"errors"
	"fmt"
	"log"
	"main/internal/dismissal"
	"main/internal/issue"
	"main/internal/object"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/savioxavier/termlink"
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

func (h *TaskHandler) AssignAllDeactivateInsightIssuesToMe() error {
	// jql := ""
	jql := `project = SD and assignee = empty and summary ~ "Деактивировать в Insight" and resolution = unresolved and "Postpone until" < endOfWeek()`

	insightIssues, err := h.issueService.GetAll(jql)
	if err != nil {
		return fmt.Errorf("failed to get all insight issues to assign: %w", err)
	}

	for _, insightIssue := range insightIssues {
		_, err = h.issueService.AssignIssueToMe(&insightIssue)
		fmt.Printf("assigning [%v] %v to self\n", insightIssue.Key, insightIssue.Fields.Summary)
		if err != nil {
			return fmt.Errorf("failed to assign issue to me: %w", err)
		}

		time.Sleep(time.Second)

	}

	return nil
}

func (h *TaskHandler) DeactivateInsight() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Деактивировать в Insight"`
	var unresolved []*jira.Issue

	deactivationIssues, err := h.issueService.GetAll(jql)
	if err != nil {
		return err
	}

	if len(deactivationIssues) == 0 {
		return errors.New("no deactivation issues")
	}

	for _, di := range deactivationIssues {
		var commentText string

		fmt.Print("found an issue: ")
		h.issueService.PrintIssue(&di)
		parentIssue, err := h.issueService.GetByID(di.Fields.Parent.ID)
		if err != nil {
			return err
		}

		fmt.Print("found a parent issue: ")
		h.issueService.PrintIssue(parentIssue)

		deactivationSubtask, err := h.issueService.GetSubtaskByComponent(parentIssue, "Insight")
		if err != nil {
			return fmt.Errorf("failed to get subtask by component: %w", err)
		}

		deactivationIssue, err := h.issueService.GetByID(deactivationSubtask.ID)
		if err != nil {
			return fmt.Errorf("failed to get deactivation issue by ID: %w", err)
		}

		subShipment, err := h.issueService.GetSubtaskByComponent(parentIssue, "Возврат оборудования")
		if err != nil {
			panic(err)
		}

		if subShipment.Fields.Status.Name != "Closed" {
			fmt.Printf("parent issue [%v] has an incomplete subshipment task\n", parentIssue.Key)
			//block deactivation issue by the aforementioned subshipment

			blockingIssue, err := h.issueService.GetByID(subShipment.ID)
			if err != nil {
				return err
			}

			fmt.Printf("blocking [%v] by [%v]\n", di.Key, blockingIssue.Key)
			_, err = h.issueService.BlockByIssue(deactivationIssue, blockingIssue)
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

		fmt.Printf("found user email %v\n", userEmail)

		user, err := h.objectService.GetUserByEmail(userEmail)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		if user == nil {
			fmt.Printf("couldn't find insight user %v\n", userEmail)
			commentText = "Пользователя нет в Insight"
		} else {
			getUserLaptopsRes, err := h.objectService.GetUserLaptops(user)
			if err != nil {
				return err
			}

			laptops := getUserLaptopsRes.ObjectEntries

			fmt.Printf("user %v has %v laptops\n", userEmail, len(laptops))

			var category string

			if len(laptops) > 0 {
				fmt.Printf("user %v still has attached laptops\n", userEmail)
				unresolved = append(unresolved, &di)
				continue
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

			commentText = "[https://wiki.sbmt.io/x/sPjivQ]"
		}

		_, err = h.issueService.Close(deactivationIssue)
		if err != nil {
			return fmt.Errorf("failed to close deactivation issue: %w", err)
		}

		fmt.Printf("adding internal comment to [%v]\n", deactivationIssue.Key)
		_, err = h.issueService.WriteInternalComment(deactivationIssue, commentText)
		if err != nil {
			return fmt.Errorf("failed to write comment: %w", err)
		}

		fmt.Println("timeout 3 sec")
		time.Sleep(3 * time.Second)

	}

	if len(unresolved) > 0 {
		fmt.Print("unresolved issues:\n")
		for i, ui := range unresolved {
			fmt.Printf("%v. ", i+1)
			h.issueService.PrintIssue(ui)
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
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
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
				return fmt.Errorf("failed to create template %v: %w", templateName, err)
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

func (h *TaskHandler) ShowIssuesWithEmptyComponent() error {
	jql := `project = SD AND component = EMPTY AND assignee in (EMPTY) AND resolution = Unresolved and updated > startOfDay()`
	for {
		fmt.Print("\033[H\033[2J")
		issues, err := h.issueService.GetAll(jql)
		if err != nil {
			return fmt.Errorf("failed to get all issues with empty component: %w", err)
		}

		for _, issue := range issues {
			summary := h.issueService.Summarize(&issue)
			issueLink := fmt.Sprintf("https://jira.sbmt.io/browse/%v", issue.Key)

			fmt.Println(termlink.Link(summary, issueLink))
		}

		time.Sleep(5 * time.Second)

	}

}

func (h *TaskHandler) UpdateBlockTraineeIssue() error {
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)

	issue, err := h.issueService.GetByID(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue by key: %w", err)
	}

	var causingIssue *jira.Issue

	issueLinks := issue.Fields.IssueLinks
	for _, issueLink := range issueLinks {
		if issueLink.Type.Inward == "is caused by" {
			causingIssue, err = h.issueService.GetByID(issueLink.InwardIssue.ID)
			if err != nil {
				return fmt.Errorf("failed to get issue %v by id: %w", issueLink.InwardIssue.Key, err)
			}
		}
	}

	email := causingIssue.Fields.Unknowns["customfield_10356"].(string)
	fmt.Printf("user email: %v\n", email)

	for _, st := range issue.Fields.Subtasks {
		subtaskIssue, err := h.issueService.GetByID(st.ID)
		if err != nil {
			return fmt.Errorf("failed to get subtask %v: %w", st.Key, err)
		}
		h.issueService.PrintIssue(subtaskIssue)

		currentSummary := strings.TrimSpace(subtaskIssue.Fields.Summary)

		newSummary := currentSummary + " " + email

		fmt.Printf("new summary: %v\n", newSummary)

		type Fields struct {
			Summary string `json:"summary" structs:"summary"`
		}

		c := map[string]interface{}{
			"fields": Fields{
				Summary: newSummary,
			},
		}

		_, err = h.issueService.Update(subtaskIssue, c)
		if err != nil {
			return fmt.Errorf("failed to update summary for %v: %w", subtaskIssue.Key, err)
		}

		time.Sleep(3 * time.Second)
	}

	return nil

}
