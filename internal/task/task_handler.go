package task

import (
	"errors"
	"fmt"
	"io"
	"main/internal/ad"
	"main/internal/dismissal"
	"main/internal/issue"
	"main/internal/object"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/savioxavier/termlink"
)

const longTimeout = 2
const shortTimeout = 1

type TaskHandler struct {
	issueService     issue.IssueService
	objectService    object.ObjectService
	dismissalService dismissal.DismissalService
	adService        ad.AdService
}

func NewTaskHandler(is *issue.IssueService, os *object.ObjectService, ds *dismissal.DismissalService, as *ad.AdService) *TaskHandler {
	return &TaskHandler{issueService: *is, objectService: *os, dismissalService: *ds, adService: *as}
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

		time.Sleep(shortTimeout * time.Second)

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

			fmt.Println("timeout 5 sec")
			time.Sleep(5 * time.Second)

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

		fmt.Println("timeout 5 sec")
		time.Sleep(5 * time.Second)

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

		time.Sleep(longTimeout * time.Second)

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
	//email, _ := issue.Fields.Unknowns.Value(EMAIL_FIELD_KEY)
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

		time.Sleep(longTimeout * time.Second)
	}

	return nil

}

func (h *TaskHandler) RemoveVPNGroupsFromUsers() error {
	f, err := os.OpenFile("vpn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(f, os.Stdout)
	log.SetOutput(wrt)

	var failedUsers []string
	failedFile, err := os.Create("failed_users.txt")
	if err != nil {
		return fmt.Errorf("failed to create failed users file: %w", err)
	}
	defer failedFile.Close()

	prefix := "res.vpn."
	usersFilename := "users.txt"

	log.Info(fmt.Sprintf("Opening %s file", usersFilename))
	file, err := os.Open(usersFilename)
	if err != nil {
		log.Error(fmt.Errorf("failed to open file %s: %w", usersFilename, err))
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Error(fmt.Errorf("failed to read file %s: %w", usersFilename, err))
		return err
	}

	emails := strings.Split(string(content), "\n")

	for _, email := range emails {
		log.Info(fmt.Sprintf("Getting user by email %s", email))
		user, err := h.adService.GetByEmail(email)
		if err != nil {
			log.Error(fmt.Errorf("failed to get user %s by email: %w", email, err))
			continue
		}

		userGroups := h.adService.GetUserGroups(user)

		for _, groupDN := range userGroups {
			groupCN, err := h.adService.ExtractCNFromDN(groupDN)
			if err != nil {
				log.Error(fmt.Errorf("failed to extract group CN from DN %s: %w", groupDN, err))
				continue
			}

			if strings.HasPrefix(groupCN, prefix) {
				group := &ldap.Entry{DN: groupDN}

				log.Info(fmt.Sprintf("Removing group %s from user %s", groupCN, user.GetAttributeValue("sAMAccountName")))

				_, err := h.adService.RemoveUserFromGroup(user, group)
				if err != nil {
					log.Error(fmt.Errorf("failed to remove user %s from group %s: %w", email, groupCN, err))
					failedUsers = append(failedUsers, email)
					continue
				}
				time.Sleep(time.Duration(shortTimeout) * time.Second)
			}

		}
	}

	for _, failedUser := range failedUsers {
		_, err := failedFile.WriteString(failedUser + "\n")
		if err != nil {
			log.Error(fmt.Errorf("failed to write failed user %s to file: %v", failedUser, err))
		}
	}

	log.Info("Task completed successfully")
	return nil
}

func (h *TaskHandler) AddUserToGroupFromJiraIssue() error {
	var adGroupCN string
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)
	// issueKey := "SD-735229"

	log.Info(fmt.Sprintf("Getting issue by key %v", issueKey))
	issue, err := h.issueService.GetByID(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue by key: %w", err)
	}
	log.Info(fmt.Sprintf("Found issue: %v", h.issueService.Summarize(issue)))

	roleInfo := issue.Fields.Unknowns["customfield_13063"].([]interface{})[0].(string) //unreliable
	roleInfoArray := strings.Split(roleInfo, " ")
	informationResourceKeyRaw := roleInfoArray[len(roleInfoArray)-1]
	informationResourceKey := informationResourceKeyRaw[1 : len(informationResourceKeyRaw)-1]

	log.Info(fmt.Sprintf("Getting information resource by key %v", informationResourceKey))

	informationResource, err := h.objectService.GetByISC(informationResourceKey)
	if err != nil {
		return fmt.Errorf("failed to get information resource by key: %w", err)
	}
	log.Info(fmt.Sprintf("Found information resource: %v", informationResource.Label))

	for _, attribute := range informationResource.Attributes {
		if attribute.ObjectTypeAttributeID == 8527 {
			adGroupCN = attribute.ObjectAttributeValues[0].Value
		}
	}

	if adGroupCN == "" {
		return errors.New("empty ad group CN")
	}
	log.Info(fmt.Sprintf("Found AD group: %v", adGroupCN))

	summary := strings.Split(issue.Fields.Summary, " ")
	email := summary[len(summary)-1]

	log.Info(fmt.Sprintf("Getting AD user by email: %v", email))
	user, err := h.adService.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}
	group, err := h.adService.GetByCN(adGroupCN)
	if err != nil {
		return fmt.Errorf("failed to get group by cn: %w", err)
	}

	log.Info(fmt.Sprintf("Adding user %v to group %v", user.GetAttributeValue("mail"), group.GetAttributeValue("cn")))
	_, err = h.adService.AddUserToGroup(user, group)
	if err != nil {
		return fmt.Errorf("failed to add user %v to group %v : %w", user.GetAttributeValue("mail"), group.GetAttributeValue("cn"), err)
	}

	commentText := "[https://wiki.sbmt.io/x/WcPivQ]"

	log.Info(fmt.Sprintf("adding internal comment to [%v]\n", issue.Key))
	_, err = h.issueService.WriteInternalComment(issue, commentText)
	if err != nil {
		return fmt.Errorf("failed to write comment: %w", err)
	}

	return nil
}

func (h *TaskHandler) CheckUserStatus() error {
	//project = sd and assignee = empty and resolution = unresolved and text ~ "Блокировка УЗ в AD" and due < endOfWeek() and type = Sub-task and component = AD
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)
	// issueKey := "SD-731877"

	log.Info(fmt.Sprintf("Getting issue by key %v", issueKey))
	issue, err := h.issueService.GetByID(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue by key: %w", err)
	}
	log.Info(fmt.Sprintf("Found issue: %v", h.issueService.Summarize(issue)))
	//summary: 19.06 Создать УЗ AD для заявки на стажера Иванов Иван Иванович
	summary := strings.Split(issue.Fields.Summary, " ") // won't work
	email := summary[len(summary)-1]

	log.Info(fmt.Sprintf("Getting AD user by email: %v", email))
	//! what if more than one user has the same name
	user, err := h.adService.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	flags := user.GetAttributeValue("userAccountControl") //512 -- active, 514 -- inactive

	flagsInt, err := strconv.Atoi(flags)
	if err != nil {
		return fmt.Errorf("failed to convert userAccountControl to int")
	}

	status := "unknown"

	if flagsInt == 512 {
		status = "active"
		_, err := h.issueService.BlockUntilTomorrow(issue)
		if err != nil {
			return fmt.Errorf("failed to block issue until tomorrow: %w", err)
		}
		//block until tomorrow
	}

	if flagsInt == 514 {
		//add comment, close the issue
		status = "inactive"
		_, err := h.issueService.Close(issue)
		if err != nil {
			return fmt.Errorf("failed to close issue: %w", err)
		}

		_, err = h.issueService.WriteInternalComment(issue, "Заблокировал idm")
		if err != nil {
			return fmt.Errorf("failed to write internal comment: %w", err)
		}

		_, err = h.issueService.WriteInternalComment(issue, "[https://wiki.sbmt.io/x/mB6zsg]")
		if err != nil {
			return fmt.Errorf("failed to write internal comment: %w", err)
		}
	}

	log.Info(fmt.Sprintf("user's %v status is %v", email, status))

	return nil

}

func (h *TaskHandler) MoveUsersToNewOU() error {
	f, err := os.OpenFile("ou.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(f, os.Stdout)
	log.SetOutput(wrt)

	usersFilename := "users.txt"

	log.Info(fmt.Sprintf("Opening %s file", usersFilename))
	file, err := os.Open(usersFilename)
	if err != nil {
		log.Error(fmt.Errorf("failed to open file %s: %w", usersFilename, err))
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Error(fmt.Errorf("failed to read file %s: %w", usersFilename, err))
		return err
	}

	commonNames := strings.Split(string(content), "\n")

	for _, cn := range commonNames {
		user, err := h.adService.GetByCN(cn)
		if err != nil {
			log.Error(fmt.Errorf("failed to get user %v by cn", cn))
			return err
		}

		newSup := "OU=External,OU=Users,OU=SBMT,DC=sbermarket,DC=ru"
		// newDN := fmt.Sprintf("CN=%v,OU=E-comm,OU=Users,OU=SBMT,DC=sbermarket,DC=ru", cn)

		log.Info(fmt.Sprintf("Moving user %v to %v", cn, newSup))
		_, err = h.adService.UpdateDN(user, newSup)
		if err != nil {
			log.Error(fmt.Errorf("failed to move user %v to %v", cn, newSup))
			return err
		}

		time.Sleep(shortTimeout * time.Second)

	}

	return nil
}
