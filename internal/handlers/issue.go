package handlers

import (
	"errors"
	"fmt"
	"main/internal/interfaces"
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

func NewIssueHandler(issueService interfaces.IssueService, activeDirectoryService interfaces.ActiveDirectoryService, assetService interfaces.AssetService) *issueHandler {
	return &issueHandler{issueService: issueService, activeDirectoryService: activeDirectoryService, assetService: assetService}
}

func (issueHandler *issueHandler) DeactivateInsight() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = "Open" AND summary ~ "Деактивировать в Insight"`
	var unresolvedIssues []*jira.Issue

	deactivateInsightIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return err
	}

	if len(deactivateInsightIssues) == 0 {
		return errors.New("no deactivation issues")
	}

	for _, deactivateInsightIssue := range deactivateInsightIssues {
		var commentText string

		log.Info(fmt.Sprintf("getting parent issue for: %v", issueHandler.issueService.Summarize(&deactivateInsightIssue)))
		parentIssue, err := issueHandler.issueService.GetByID(deactivateInsightIssue.Fields.Parent.ID)
		if err != nil {
			return err
		}

		log.Info(fmt.Sprintf("getting return equipment subtask for: %v", parentIssue.Key))
		returnEquipmentSubtask, err := issueHandler.issueService.GetSubtaskByComponent(parentIssue, "Возврат оборудования")
		if err != nil {
			panic(err)
		}

		if returnEquipmentSubtask.Fields.Status.Name != "Closed" {
			log.Info(fmt.Sprintf("parent issue %v has incomplete return equipment task", parentIssue.Key))

			blockingIssue, err := issueHandler.issueService.GetByID(returnEquipmentSubtask.ID)
			if err != nil {
				return err
			}

			log.Info(fmt.Sprintf("blocking %v by %v", issueHandler.issueService.Summarize(&deactivateInsightIssue), issueHandler.issueService.Summarize(blockingIssue)))
			_, err = issueHandler.issueService.BlockByIssue(&deactivateInsightIssue, blockingIssue)
			if err != nil {
				return err
			}

			log.Info("timeout 5 sec")
			time.Sleep(5 * time.Second)

			continue
		}

		summaryFields := strings.Fields(deactivateInsightIssue.Fields.Summary)
		userEmail := summaryFields[len(summaryFields)-1]

		log.Info(fmt.Sprintf("found user email %v", userEmail))

		user, err := issueHandler.assetService.GetUserByEmail(userEmail)
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		if user == nil {
			log.Info(fmt.Sprintf("couldn't find insight user %v", userEmail))
			commentText = "Пользователя нет в Insight"
		} else {
			getUserLaptopsRes, err := issueHandler.assetService.GetUserLaptops(user)
			if err != nil {
				return err
			}

			laptops := getUserLaptopsRes.ObjectEntries

			log.Info(fmt.Sprintf("user %v has %v laptops", userEmail, len(laptops)))

			if len(laptops) > 0 {
				log.Info(fmt.Sprintf("user %v still has attached laptops", userEmail))
				unresolvedIssues = append(unresolvedIssues, &deactivateInsightIssue)
				continue
			}

			category := "BYOD"

			log.Info(fmt.Sprintf("changing %v's status to %v", user.ObjectKey, category))
			_, err = issueHandler.assetService.SetUserCategory(user, category)
			if err != nil {
				return fmt.Errorf("failed to set user category: %w", err)
			}

			log.Info(fmt.Sprintf("disabling %v", user.ObjectKey))
			_, err = issueHandler.assetService.DisableUser(user)
			if err != nil {
				return fmt.Errorf("failed to disable user: %w", err)
			}

			commentText = "[https://wiki.sbmt.io/x/sPjivQ]"
		}

		_, err = issueHandler.issueService.Close(&deactivateInsightIssue)
		if err != nil {
			return fmt.Errorf("failed to close deactivation issue: %w", err)
		}

		log.Info(fmt.Sprintf("adding internal comment to %v", deactivateInsightIssue.Key))
		_, err = issueHandler.issueService.WriteInternalComment(&deactivateInsightIssue, commentText)
		if err != nil {
			return fmt.Errorf("failed to write internal comment: %w", err)
		}

		log.Info("timeout 5 sec")
		time.Sleep(5 * time.Second)

	}

	if len(unresolvedIssues) > 0 {
		log.Info("unresolved issues")
		for i, unresolvedIssue := range unresolvedIssues {
			log.Info(fmt.Sprintf("%v. %v", i+1, issueHandler.issueService.Summarize(unresolvedIssue)))
		}
	}

	return nil
}

func (issueHandler *issueHandler) AddUserToGroupFromJiraIssue() error {
	var adGroupCN string
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)

	log.Info(fmt.Sprintf("Getting issue by key %v", issueKey))
	issue, err := issueHandler.issueService.GetByID(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue by key: %w", err)
	}
	log.Info(fmt.Sprintf("Found issue: %v", issueHandler.issueService.Summarize(issue)))

	roleInfo := issue.Fields.Unknowns["customfield_13063"].([]interface{})[0].(string) //unreliable
	roleInfoArray := strings.Split(roleInfo, " ")
	informationResourceKeyRaw := roleInfoArray[len(roleInfoArray)-1]
	informationResourceKey := informationResourceKeyRaw[1 : len(informationResourceKeyRaw)-1]

	log.Info(fmt.Sprintf("Getting information resource by key %v", informationResourceKey))

	informationResource, err := issueHandler.assetService.GetByISC(informationResourceKey)
	if err != nil {
		return fmt.Errorf("failed to get information resource by key: %w", err)
	}
	log.Info(fmt.Sprintf("Found information resource: %v", informationResource.Label))

	for _, attribute := range informationResource.Attributes {
		if attribute.ObjectTypeAttributeID == 8527 {
			adGroupCN = strings.TrimSpace(attribute.ObjectAttributeValues[0].Value)
		}
	}

	if adGroupCN == "" {
		return errors.New("empty ad group CN")
	}

	fmt.Println(adGroupCN)

	group, err := issueHandler.activeDirectoryService.GetByCN(adGroupCN)
	if err != nil {
		return fmt.Errorf("failed to get group by cn: %w", err)
	}

	log.Info(fmt.Sprintf("Found AD group: %v", adGroupCN))

	summary := strings.Split(strings.TrimSpace(issue.Fields.Summary), " ") // won't work
	email := summary[len(summary)-1]

	log.Info(fmt.Sprintf("Getting AD user by email: %v", email))
	user, err := issueHandler.activeDirectoryService.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	log.Info(fmt.Sprintf("Adding user %v to group %v", user.GetAttributeValue("mail"), group.GetAttributeValue("cn")))
	_, err = issueHandler.activeDirectoryService.AddUserToGroup(user, group)
	if err != nil {
		return fmt.Errorf("failed to add user %v to group %v : %w", user.GetAttributeValue("mail"), group.GetAttributeValue("cn"), err)
	}

	commentText := "[https://wiki.sbmt.io/x/WcPivQ]"

	log.Info(fmt.Sprintf("adding internal comment to [%v]", issue.Key))
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
	jql := `project = SD and assignee = empty and (summary ~ "Деактивировать в Insight" or summary ~ "Блокировка УЗ в AD") and component in (Insight, AD) and resolution = unresolved and "Postpone until" < endOfMonth()`

	insightIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return fmt.Errorf("failed to get all insight issues to assign: %w", err)
	}

	for _, insightIssue := range insightIssues {
		_, err = issueHandler.issueService.AssignIssueToMe(&insightIssue)
		log.Info(fmt.Sprintf("assigning [%v] %v to self\n", insightIssue.Key, insightIssue.Fields.Summary))
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

	var users []*ldap.Entry
	var groups []*ldap.Entry

	fmt.Print("enter user(s) email: ")
	fmt.Scanln(&emailList)

	fmt.Print("enter group(s) cn: ")
	fmt.Scanln(&activeDirectoryGroupCNList)

	emails := strings.Fields(emailList)
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
		log.Info("getting AD group by cn: %v", groupCN)
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

	for _, issue := range issues {
		log.Info(fmt.Sprintf("processing issue %v", issueHandler.issueService.Summarize(&issue)))
		log.Info(fmt.Sprintf("getting unresolved subtask for %v", issue.Key))
		unresolvedSubtask, err := issueHandler.issueService.GetUnresolvedSubtask(&issue)
		if err != nil {
			return fmt.Errorf("failed to get unresolved subtask: %w", err)
		}

		if unresolvedSubtask != nil {
			log.Info(fmt.Sprintf("found unresolved subtask: %v %v", unresolvedSubtask.Key, unresolvedSubtask.Fields.Summary))
			log.Info(fmt.Sprintf("blocking main issue %v by unresolved subtask %v", issue.Key, unresolvedSubtask.Key))
			_, err := issueHandler.issueService.BlockByIssue(&issue, unresolvedSubtask)
			if err != nil {
				return fmt.Errorf("failed to block by issue")
			}

			continue
		}

		log.Info("all subtasks are resolved")
		log.Info("closing the issue")
		_, err = issueHandler.issueService.Close(&issue)
		if err != nil {
			return fmt.Errorf("failed to close issue %v: %w", issue.Key, err)
		}

		var comment string

		for _, component := range issue.Fields.Components {
			if component.Name == "Dismissal" {
				comment = "[https://wiki.sbmt.io/x/jgeLvg]"
			}

			if component.Name == "Hiring" {
				comment = "[https://wiki.sbmt.io/x/ogeLvg]"
			}
		}

		log.Info("adding internal comment")

		_, err = issueHandler.issueService.WriteInternalComment(&issue, comment)
		if err != nil {
			return fmt.Errorf("failed to write internal comment to %v: %w", issue.Key, err)
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (issueHandler *issueHandler) CheckUserStatus() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Блокировка УЗ в AD для"`

	blockIssues, err := issueHandler.issueService.GetAll(jql)
	if err != nil {
		return err
	}

	if len(blockIssues) == 0 {
		return errors.New("no deactivation issues")
	}

	for _, issue := range blockIssues {
		log.Info(fmt.Sprintf("Found issue: %v", issueHandler.issueService.Summarize(&issue)))
		//summary: 19.06 Создать УЗ AD для заявки на стажера Иванов Иван Иванович
		summary := strings.Split(strings.TrimSpace(issue.Fields.Summary), " ") // won't work
		email := summary[len(summary)-1]

		log.Info(fmt.Sprintf("Getting AD user by email: %v", email))
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
			log.Info(fmt.Sprintf("user %v is active", email))
			log.Info(fmt.Sprintf("blocking issue %v until tomorrow", issue.Key))

			_, err := issueHandler.issueService.BlockUntilTomorrow(&issue)
			if err != nil {
				return fmt.Errorf("failed to block issue until tomorrow: %w", err)
			}

		}

		if flagInt == 514 {
			log.Info(fmt.Sprintf("user %v is inactive", email))

			log.Info(fmt.Sprintf("closing issue %v", issue.Key))
			_, err := issueHandler.issueService.Close(&issue)
			if err != nil {
				return fmt.Errorf("failed to close issue: %w", err)
			}

			log.Info(fmt.Sprintf("writing internal comment to %v", issue.Key))
			_, err = issueHandler.issueService.WriteInternalComment(&issue, "Заблокировал idm")
			if err != nil {
				return fmt.Errorf("failed to write internal comment: %w", err)
			}

			_, err = issueHandler.issueService.WriteInternalComment(&issue, "[https://wiki.sbmt.io/x/mB6zsg]")
			if err != nil {
				return fmt.Errorf("failed to write internal comment: %w", err)
			}
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}
