package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"main/internal/interfaces"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap/v3"
	"github.com/savioxavier/termlink"
)

type issueHandler struct {
	issueService           interfaces.IssueService
	activeDirectoryService interfaces.ActiveDirectoryService
	assetService           interfaces.AssetService
}

func NewIssueHandler(issueService interfaces.IssueService, activeDirectoryService interfaces.ActiveDirectoryService, assetService interfaces.AssetService) interfaces.IssueHandler {
	return &issueHandler{issueService: issueService, activeDirectoryService: activeDirectoryService, assetService: assetService}
}

func (issueHandler *issueHandler) ProcessDeactivateInsightAccountIssue() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = "Open" AND summary ~ "Деактивировать в Insight"`

	deactivateInsightIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return err
	}

	if len(deactivateInsightIssues) == 0 {
		log.Info("no issues found")
	}

	for _, deactivateInsightIssue := range deactivateInsightIssues {
		log.SetPrefix(deactivateInsightIssue.Key)

		log.Info(fmt.Sprintf("start processing the issue %v", issueHandler.issueService.Summarize(&deactivateInsightIssue)))

		log.Info("getting parent issue")
		parentIssue, err := issueHandler.issueService.GetByID(deactivateInsightIssue.Fields.Parent.ID)
		if err != nil {
			return fmt.Errorf("failed to get parent issue %v: %w", deactivateInsightIssue.Fields.Parent.Key, err)
		}
		log.Info(fmt.Sprintf("found parent issue: %v", issueHandler.issueService.Summarize(parentIssue)))

		componentName := "Возврат оборудования"

		log.Info(fmt.Sprintf("getting subtask by component %v", componentName))
		returnEquipmentSubtask, err := issueHandler.issueService.GetSubtaskByComponent(parentIssue, componentName)
		if err != nil {
			return fmt.Errorf("failed to get subtask by component %v: %w", componentName, err)
		}

		if returnEquipmentSubtask.Fields.Status.Name != "Closed" {
			log.Info(fmt.Sprintf("parent issue %v has incomplete return equipment task", parentIssue.Key))

			blockingIssue, err := issueHandler.issueService.GetByID(returnEquipmentSubtask.ID)
			if err != nil {
				return err
			}

			log.Info(fmt.Sprintf("blocking the issue by %v", issueHandler.issueService.Summarize(blockingIssue)))
			_, err = issueHandler.issueService.BlockByIssue(&deactivateInsightIssue, blockingIssue)
			if err != nil {
				return err
			}

			log.Info("finished processing the issue\n")
			time.Sleep(5 * time.Second)

			continue
		}

		var commentText string

		summaryFields := strings.Fields(deactivateInsightIssue.Fields.Summary)
		email := summaryFields[len(summaryFields)-1]

		log.Info(fmt.Sprintf("getting user by email %v", email))
		user, err := issueHandler.assetService.GetUserByEmail(email)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		if user != nil {
			log.Info("getting user laptops")
			getUserLaptopsRes, err := issueHandler.assetService.GetUserLaptops(user)
			if err != nil {
				return fmt.Errorf("failed to get user %v laptops: %w", email, err)
			}

			laptops := getUserLaptopsRes.ObjectEntries

			log.Info(fmt.Sprintf("user has %v attached laptops", len(laptops)))

			if len(laptops) > 0 {
				log.Info("user still has attached laptops")
				log.Info("skipping the issue...\n")

				time.Sleep(5 * time.Second)
				continue
			}

			category := "BYOD"

			log.Info(fmt.Sprintf("changing user status to %v", category))
			_, err = issueHandler.assetService.SetUserCategory(user, category)
			if err != nil {
				return fmt.Errorf("failed to set user category: %w", err)
			}

			log.Info("disabling user")
			_, err = issueHandler.assetService.DisableUser(user)
			if err != nil {
				return fmt.Errorf("failed to disable user: %w", err)
			}

			commentText = "[https://wiki.sbmt.io/x/sPjivQ]"

		} else {
			log.Info(fmt.Sprintf("couldn't find insight user %v", email))
			commentText = "Пользователя нет в Insight"
		}

		log.Info(fmt.Sprintf("adding internal comment \"%v\"", commentText))
		_, err = issueHandler.issueService.WriteInternalComment(&deactivateInsightIssue, commentText)
		if err != nil {
			return fmt.Errorf("failed to write internal comment: %w", err)
		}

		log.Info("closing the issue")
		_, err = issueHandler.issueService.Close(&deactivateInsightIssue)
		if err != nil {
			return fmt.Errorf("failed to close deactivation issue: %w", err)
		}

		log.Info("finished processing the issue\n")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (issueHandler *issueHandler) ProcessGrantAccessIssue() error {
	var activeDirectoryGroupCN string
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)

	log.SetPrefix(issueKey)
	log.Info("getting issue by key")
	issue, err := issueHandler.issueService.GetByID(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue by key: %w", err)
	}
	log.Info(fmt.Sprintf("found issue: %v", issueHandler.issueService.Summarize(issue)))
	roleInfoFieldID := "customfield_13063"

	roleField, err := issueHandler.issueService.GetCustomFieldValue(issue, roleInfoFieldID)
	if err != nil {
		return fmt.Errorf("failed to get custom field %v value", roleInfoFieldID)
	}

	roleInfo := roleField.([]interface{})[0].(string)

	log.Info(fmt.Sprintf("role info: %v", roleInfo))

	informationResourceKey, err := issueHandler.assetService.ExtractInformationResourceIdentifier(roleInfo)
	if err != nil {
		return fmt.Errorf("failed to extract information resource identifier from custom field value: %w", err)
	}

	log.Info(fmt.Sprintf("getting information resource by key %v", informationResourceKey))
	informationResource, err := issueHandler.assetService.GetByISC(informationResourceKey)
	if err != nil {
		return fmt.Errorf("failed to get information resource by key: %w", err)
	}
	log.Info(fmt.Sprintf("found information resource: %v", informationResource.Label))

	for _, attribute := range informationResource.Attributes {
		if attribute.ObjectTypeAttributeID == 8527 {
			activeDirectoryGroupCN = strings.TrimSpace(attribute.ObjectAttributeValues[0].Value)
		}
	}

	if activeDirectoryGroupCN == "" {
		return errors.New("empty ad group CN")
	}

	log.Info(fmt.Sprintf("getting active directory group by CN: %v", activeDirectoryGroupCN))
	group, err := issueHandler.activeDirectoryService.GetByCN(activeDirectoryGroupCN)
	if err != nil {
		return fmt.Errorf("failed to get group by cn: %w", err)
	}

	summary := strings.Fields(issue.Fields.Summary)
	email := summary[len(summary)-1]

	log.Info(fmt.Sprintf("getting active directory user by email: %v", email))
	user, err := issueHandler.activeDirectoryService.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	log.Info("adding user to group")
	_, err = issueHandler.activeDirectoryService.AddUserToGroup(user, group)
	if err != nil {
		return fmt.Errorf("failed to add user %v to group %v : %w", user.GetAttributeValue("mail"), group.GetAttributeValue("cn"), err)
	}

	commentText := "[https://wiki.sbmt.io/x/WcPivQ]"

	log.Info(fmt.Sprintf("adding internal comment \"%v\"", commentText))
	_, err = issueHandler.issueService.WriteInternalComment(issue, commentText)
	if err != nil {
		return fmt.Errorf("failed to write comment: %w", err)
	}

	return nil
}

func (issueHandler *issueHandler) UpdateBlockTraineeIssue() error {
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)

	issue, err := issueHandler.issueService.GetByID(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue by key: %w", err)
	}

	var causingIssue *jira.Issue

	issueLinks := issue.Fields.IssueLinks
	for _, issueLink := range issueLinks {
		if issueLink.Type.Inward == "is caused by" {
			causingIssue, err = issueHandler.issueService.GetByID(issueLink.InwardIssue.ID)
			if err != nil {
				return fmt.Errorf("failed to get issue %v by id: %w", issueLink.InwardIssue.Key, err)
			}
		}
	}

	email := causingIssue.Fields.Unknowns["customfield_10356"].(string)
	//email, _ := issue.Fields.Unknowns.Value(EMAIL_FIELD_KEY)
	fmt.Printf("user email: %v\n", email)

	for _, st := range issue.Fields.Subtasks {
		subtaskIssue, err := issueHandler.issueService.GetByID(st.ID)
		if err != nil {
			return fmt.Errorf("failed to get subtask %v: %w", st.Key, err)
		}
		issueHandler.issueService.PrintIssue(subtaskIssue)

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

		_, err = issueHandler.issueService.Update(subtaskIssue, c)
		if err != nil {
			return fmt.Errorf("failed to update summary for %v: %w", subtaskIssue.Key, err)
		}

		time.Sleep(time.Second)
	}

	return nil

}

func (issueHandler *issueHandler) ShowIssuesWithEmptyComponent() error {
	jql := `project = SD AND component = EMPTY AND assignee in (EMPTY) AND resolution = Unresolved and updated > startOfDay()`
	for {
		fmt.Print("\033[H\033[2J")
		issues, err := issueHandler.issueService.GetAll(jql)
		if err != nil {
			return fmt.Errorf("failed to get all issues with empty component: %w", err)
		}

		for _, issue := range issues {
			summary := issueHandler.issueService.Summarize(&issue)
			issueLink := fmt.Sprintf("https://jira.sbmt.io/browse/%v", issue.Key)

			fmt.Println(termlink.Link(summary, issueLink))
		}

		time.Sleep(5 * time.Second)

	}

}

func (issueHandler *issueHandler) AssignAllDeactivateInsightIssuesToMe() error {
	jql := `project = SD and assignee = empty and (summary ~ "Деактивировать в Insight" or summary ~ "Блокировка УЗ в AD") and component in (Insight, AD) and resolution = unresolved and "Postpone until" < endOfWeek()`

	foundIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return fmt.Errorf("failed to get all insight issues to assign: %w", err)
	}

	if len(foundIssues) == 0 {
		log.Info("no issues found")
	}

	for _, foundIssue := range foundIssues {
		log.SetPrefix(foundIssue.Key)
		log.Info(fmt.Sprintf("assigning %v to me", foundIssue.Fields.Summary))
		_, err = issueHandler.issueService.AssignIssueToMe(&foundIssue)
		if err != nil {
			return fmt.Errorf("failed to assign issue to me: %w", err)
		}

		time.Sleep(time.Second)

	}

	return nil
}

func (issueHandler *issueHandler) AddUserToGroupFromCLI() error {
	var emailList string
	var activeDirectoryGroupCNList string
	scanner := bufio.NewScanner(os.Stdin)

	var users []*ldap.Entry
	var groups []*ldap.Entry

	fmt.Print("enter user(s) email: ")
	if scanner.Scan() {
		emailList = scanner.Text()
	}
	emails := strings.Fields(emailList)

	fmt.Print("enter group(s) cn: ")
	if scanner.Scan() {
		activeDirectoryGroupCNList = scanner.Text()
	}
	groupCNs := strings.Fields(activeDirectoryGroupCNList)

	for _, email := range emails {
		log.Info(fmt.Sprintf("getting AD user by email: %v", email))
		user, err := issueHandler.activeDirectoryService.GetByEmail(email)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		users = append(users, user)
	}

	for _, groupCN := range groupCNs {
		log.Info(fmt.Sprintf("getting AD group by cn: %v", groupCN))
		group, err := issueHandler.activeDirectoryService.GetByCN(groupCN)
		if err != nil {
			return fmt.Errorf("failed to get group by cn: %w", err)
		}

		groups = append(groups, group)
	}

	for _, user := range users {
		for _, group := range groups {
			log.Info(fmt.Sprintf("adding user %v to group %v", user.GetAttributeValue("mail"), group.GetAttributeValue("cn")))
			_, err := issueHandler.activeDirectoryService.AddUserToGroup(user, group)
			if err != nil {
				return fmt.Errorf("failed to add user %v to group %v: %w", user.GetAttributeValue("mail"), group.GetAttributeValue("cn"), err)
			}

		}
	}

	return nil
}

func (issueHandler *issueHandler) ProcessDismissalOrHiringIssue() error {
	jql := `project = "IT Support" and assignee = currentUser() and component in (Dismissal, Hiring) and (text ~ "Прием" or text ~ "Увольнение") and status = open`

	issues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return fmt.Errorf("failed to get all issues: %w", err)
	}

	if len(issues) == 0 {
		log.Info("no issues found")
	}

	for _, issue := range issues {
		log.SetPrefix(issue.Key)

		log.Info(fmt.Sprintf("start processing the issue %v", issueHandler.issueService.Summarize(&issue)))
		log.Info("getting unresolved subtask")
		unresolvedSubtask, err := issueHandler.issueService.GetUnresolvedSubtask(&issue)
		if err != nil {
			return fmt.Errorf("failed to get unresolved subtask: %w", err)
		}

		if unresolvedSubtask == nil {
			log.Info("all subtasks are resolved")

			var commentText string

			for _, component := range issue.Fields.Components {
				if component.Name == "Dismissal" {
					commentText = "[https://wiki.sbmt.io/x/jgeLvg]"
				}

				if component.Name == "Hiring" {
					commentText = "[https://wiki.sbmt.io/x/ogeLvg]"
				}
			}

			log.Info(fmt.Sprintf("adding internal comment %v", commentText))

			_, err = issueHandler.issueService.WriteInternalComment(&issue, commentText)
			if err != nil {
				return fmt.Errorf("failed to write internal comment to %v: %w", issue.Key, err)
			}

			log.Info("closing the issue")
			_, err = issueHandler.issueService.Close(&issue)
			if err != nil {
				return fmt.Errorf("failed to close issue %v: %w", issue.Key, err)
			}
		} else {
			log.Info(fmt.Sprintf("found unresolved subtask: %v %v", unresolvedSubtask.Key, unresolvedSubtask.Fields.Summary))
			log.Info(fmt.Sprintf("blocking main issue by unresolved subtask %v", unresolvedSubtask.Key))
			_, err := issueHandler.issueService.BlockByIssue(&issue, unresolvedSubtask)
			if err != nil {
				return fmt.Errorf("failed to block by issue")
			}
		}

		log.Info("finished processing the issue\n")
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (issueHandler *issueHandler) ProcessDisableActiveDirectoryAccountIssue() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Блокировка УЗ в AD для"`

	blockIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return err
	}

	if len(blockIssues) == 0 {
		log.Info("no issues found")
	}

	for _, issue := range blockIssues {
		log.SetPrefix(issue.Key)
		log.Info(fmt.Sprintf("start processing the issue %v", issueHandler.issueService.Summarize(&issue)))

		summary := strings.Split(strings.TrimSpace(issue.Fields.Summary), " ")
		email := summary[len(summary)-1]

		log.Info(fmt.Sprintf("getting active directory user by email: %v", email))
		user, err := issueHandler.activeDirectoryService.GetByEmail(email)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		flag := user.GetAttributeValue("userAccountControl") //512 -- active, 514 -- inactive

		flagInt, err := strconv.Atoi(flag)
		if err != nil {
			return fmt.Errorf("failed to convert userAccountControl to int")
		}

		if flagInt == 512 {
			log.Info("user is active")
			log.Info("blocking the issue until tomorrow")
			_, err := issueHandler.issueService.BlockUntilTomorrow(&issue)
			if err != nil {
				return fmt.Errorf("failed to block issue until tomorrow: %w", err)
			}

		}

		if flagInt == 514 {
			log.Info("user is inactive")
			log.Info("closing the issue")
			_, err := issueHandler.issueService.Close(&issue)
			if err != nil {
				return fmt.Errorf("failed to close issue: %w", err)
			}

			commentText := "Заблокировал idm"
			log.Info(fmt.Sprintf("writing internal comment \"%v\"", commentText))
			_, err = issueHandler.issueService.WriteInternalComment(&issue, commentText)
			if err != nil {
				return fmt.Errorf("failed to write internal comment: %w", err)
			}

			commentText = "[https://wiki.sbmt.io/x/mB6zsg]"
			log.Info(fmt.Sprintf("writing internal comment \"%v\"", commentText))
			_, err = issueHandler.issueService.WriteInternalComment(&issue, commentText)
			if err != nil {
				return fmt.Errorf("failed to write internal comment: %w", err)
			}
		}

		log.Info("finished processing the issue\n")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (issueHandler *issueHandler) ProcessReturnCCEquipmentIssue() error {
	jql := `project = "IT Support" and assignee = currentUser() and status = Analysis  and component = "Возврат оборудования" and summary ~ "Проверить наличие техники и организовать ее забор для" and Asset = empty`

	jobTitleCC := "оператор"

	returnEquipmentIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return fmt.Errorf("failed to get issues: %w", err)
	}

	if len(returnEquipmentIssues) == 0 {
		log.Info("no issues found")
	}

	for _, returnEquipmentIssue := range returnEquipmentIssues {
		log.SetPrefix(returnEquipmentIssue.Key)
		log.Info(fmt.Sprintf("start processing the issue %v", issueHandler.issueService.Summarize(&returnEquipmentIssue)))

		summaryFields := strings.Fields(returnEquipmentIssue.Fields.Summary)
		email := summaryFields[len(summaryFields)-1]

		log.Info(fmt.Sprintf("getting asset user by email %v", email))
		user, err := issueHandler.assetService.GetUserByEmail(email)
		if err != nil {
			return fmt.Errorf("failed to get user %v by email: %w", email, err)
		}

		log.Info(fmt.Sprintf("getting %v's laptops", email))
		getUserLaptopsRes, err := issueHandler.assetService.GetUserLaptops(user)
		if err != nil {
			return fmt.Errorf("failed to get user's %v laptops: %w", email, err)
		}

		laptops := getUserLaptopsRes.ObjectEntries
		log.Info(fmt.Sprintf("user %v has %v attached laptops", email, len(laptops)))
		if len(laptops) > 0 {
			log.Info("user still has attached laptops, skipping the issue...\n")
			time.Sleep(5 * time.Second)
			continue
		}

		parentIssueID := returnEquipmentIssue.Fields.Parent.ID
		parentIssue, err := issueHandler.issueService.GetByID(parentIssueID)
		if err != nil {
			return fmt.Errorf("failed to get parent issue %v by key: %w", parentIssueID, err)
		}
		log.Info(fmt.Sprintf("found parent issue: %v", issueHandler.issueService.Summarize(parentIssue)))

		log.Info(fmt.Sprintf("getting %v's job title", email))
		jobTitleField, err := issueHandler.issueService.GetCustomFieldValue(parentIssue, "customfield_10197")
		if err != nil {
			return fmt.Errorf("failed to get custom field value for %v: %w", parentIssue.Key, err)
		}
		jobTitle := jobTitleField.(string)
		log.Info(fmt.Sprintf("%v's job title is \"%v\"", email, jobTitle))

		if strings.Trim(strings.ToLower(jobTitle), " ") != jobTitleCC {
			log.Info(fmt.Sprintf("job title \"%v\" doesn't equal to \"%v\"", jobTitle, jobTitleCC))
			log.Info("skipping the issue...\n")
			time.Sleep(5 * time.Second)
			continue
		}

		log.Info("declining the issue")
		_, err = issueHandler.issueService.Decline(&returnEquipmentIssue)
		if err != nil {
			return fmt.Errorf("failed to decline issue %v: %w", returnEquipmentIssue.Key, err)
		}

		commentText := "Оборудование не отправлялось"

		log.Info(fmt.Sprintf("writing internal comment \"%v\"", commentText))
		_, err = issueHandler.issueService.WriteInternalComment(&returnEquipmentIssue, commentText)
		if err != nil {
			return fmt.Errorf("failed to write comment to %v: %w", returnEquipmentIssue.Key, err)
		}

		log.Info("finished processing the issue\n")
		time.Sleep(5 * time.Second)

	}
	return nil
}
