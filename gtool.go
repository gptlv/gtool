package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap"
	"github.com/gocarina/gocsv"
	"github.com/goodsign/monday"
	"github.com/savioxavier/termlink"
)

type GTool interface {
	DeactivateInsight() error
	GrantAccess() error
	UpdateBlockTraineeIssue() error
	ShowEmpty() error
	AssignAll() error
	ProcessStaff() error
	ProcessAD() error
	PrintDescription() error
	GenerateRecords() error
	AddGroup() error
}

type gtool struct {
	client *jira.Client
	conn   *ldap.Conn
}

func New(client *jira.Client, conn *ldap.Conn) GTool {
	return &gtool{client: client, conn: conn}
}

func (g *gtool) DeactivateInsight() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = "Open" AND summary ~ "Деактивировать в Insight"`

	issues, _, err := g.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 1000})
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		log.Info("no issues found")
	}

	for _, issue := range issues {
		log.SetPrefix(issue.Key)
		log.Info("start processing the issue")
		err := g.processDeactivateInsightIssue(&issue)
		if err != nil {
			log.Error(fmt.Errorf("failed to process issue: %w", err))
		}
		log.Info("finished processing the issue\n")

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (g *gtool) GrantAccess() error {
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)

	log.SetPrefix(issueKey)
	log.Info("getting issue by key")
	issue, _, err := g.client.Issue.Get(issueKey, nil)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	return g.processGrantAccessIssue(issue)
}

func (g *gtool) UpdateBlockTraineeIssue() error {
	var issueKey string

	fmt.Print("enter issue key: ")
	fmt.Scanln(&issueKey)

	issue, _, err := g.client.Issue.Get(issueKey, nil)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	return g.processBlockTraineeIssue(issue)
}

func (g *gtool) ShowEmpty() error {
	jql := `project = SD AND component = EMPTY AND assignee in (EMPTY) AND resolution = Unresolved and updated > startOfDay()`

	for {
		fmt.Print("\033[H\033[2J")
		issues, _, err := g.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 1000})
		if err != nil {
			return fmt.Errorf("failed to search for issues: %w", err)
		}

		for _, issue := range issues {
			issueLink := fmt.Sprintf("https://jira.sbmt.io/browse/%s", issue.Key)
			fmt.Println(termlink.Link(fmt.Sprint(issue), issueLink))
		}

		time.Sleep(5 * time.Second)

	}

}

func (g *gtool) AssignAll() error {
	jql := `(project = SD and assignee = empty and (summary ~ "Деактивировать в Insight" or summary ~ "Блокировка УЗ в AD") and component in (Insight, AD) and resolution = unresolved and "Postpone until" < endOfMonth()) or (project = sd and (summary ~ "увольнение") and status = open and assignee = empty) or (project = sd and (summary ~ "прием") and status = open and assignee = empty)`

	issues, _, err := g.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 1000})
	if err != nil {
		return fmt.Errorf("failed to search for issues: %w", err)
	}

	if len(issues) == 0 {
		log.Info("no issues found")
	}

	user, _, err := g.client.User.GetSelf()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	payload := new(jira.User)
	payload.Name = user.Name

	for _, issue := range issues {
		log.SetPrefix(issue.Key)
		log.Info(fmt.Sprintf("assigning %s to %s", issue, user.Name))

		_, err := g.client.Issue.UpdateAssignee(issue.ID, payload)
		if err != nil {
			return fmt.Errorf("failed to update assignee: %w", err)
		}

		time.Sleep(time.Second)

	}

	return nil
}

func (g *gtool) ProcessStaff() error {
	jql := `project = "IT Support" and assignee = currentUser() and component in (Dismissal, Hiring) and (text ~ "Прием" or text ~ "Увольнение" or text ~ "Заявка на стажера") and status = open`

	issues, _, err := g.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 1000, Expand: "transitions"})
	if err != nil {
		return fmt.Errorf("failed to search for issues: %w", err)
	}

	if len(issues) == 0 {
		log.Info("no issues found")
	}

	for _, issue := range issues {
		log.SetPrefix(issue.Key)
		log.Info("start processing the issue")
		if err := g.processStaffIssue(&issue); err != nil {
			return err
		}
		log.Info("finished processing the issue\n")

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (g *gtool) ReturnEquipment() error {
	jql := `project = "IT Support" and assignee = currentUser() and status = Analysis  and component = "Возврат оборудования" and summary ~ "Проверить наличие техники и организовать ее забор для" and Asset = empty`

	issues, _, err := g.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 1000})
	if err != nil {
		return fmt.Errorf("failed to get issues: %w", err)
	}

	if len(issues) == 0 {
		log.Info("no issues found")
	}

	for _, issue := range issues {
		log.SetPrefix(issue.Key)
		err := g.processReturnEquipmentIssue(&issue)
		if err != nil {
			return fmt.Errorf("failed to process %s", issue.Key)
		}
		log.Info("finished processing the issue\n")

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (g *gtool) ProcessAD() error {
	jql := `project = "IT Support" AND assignee = currentUser() AND status = open AND summary ~ "Блокировка УЗ в AD для"`

	issues, _, err := g.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 1000})
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		log.Info("no issues found")
	}

	for _, issue := range issues {
		log.SetPrefix(issue.Key)
		log.Info("start processing the issue")
		err := g.processDisableActiveDirectoryIssue(&issue)
		if err != nil {
			return fmt.Errorf("failed to process %s: %w", issue.Key, err)
		}
		log.Info("finished processing the issue\n")

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (g *gtool) PrintDescription() error {
	var email string

	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	userList, _, err := g.getUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(userList.ObjectEntries) == 0 {
		return fmt.Errorf("no user found")
	}

	user := &userList.ObjectEntries[0]

	laptopsRes, _, err := g.getUserLaptops(user)
	if err != nil {
		return fmt.Errorf("failed to get user laptops: %w", err)
	}

	for _, laptop := range laptopsRes.ObjectEntries {
		description, err := g.getObjectDescription(&laptop)
		if err != nil {
			return fmt.Errorf("failed to get description for %s: %w", laptop.ObjectKey, err)
		}

		fmt.Print(description)
	}

	return nil
}

func (g *gtool) GenerateRecords() error {
	_, _, err := g.client.Object.Get("ISC-192756", nil)
	if err != nil {
		return err
	}
	records := []*Record{}

	inputFile := config.WriteOff.InputFile
	outputFile := config.WriteOff.OutputFile

	log.Info(fmt.Sprintf("input file name: %v", inputFile))
	log.Info(fmt.Sprintf("output file name: %v\n", outputFile))

	log.Info(fmt.Sprintf("reading csv input file %v\n", inputFile))
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	log.Info("unmarshaling input file\n")
	err = gocsv.UnmarshalFile(input, &records)
	if err != nil {
		return fmt.Errorf("failed to unmarshal input file %s: %w", inputFile, err)
	}

	t := time.Now()
	layout := "2 January 2006"
	date := monday.Format(t, layout, monday.LocaleRuRU)

	for _, record := range records {
		log.SetPrefix(record.ISC)
		log.Info("getting laptop by isc")
		laptop, _, err := g.client.Object.Get(record.ISC, nil)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		log.Info("getting laptop description")
		description, err := g.getObjectDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		record.ObjectDescription = description
		record.TeamLead = config.WriteOff.TeamLead
		record.DepartmentLead = config.WriteOff.DepartmentLead
		record.Director = config.WriteOff.Director
		record.Date = date

		log.Info("finished processing the asset\n")
		log.SetPrefix("")
	}

	log.Info("opening output file")
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	log.Info("marshaling output file")
	err = gocsv.MarshalFile(&records, output)
	if err != nil {
		return fmt.Errorf("failed to unmarshal csv file: %w", err)
	}

	log.Info("finished generating write-off records\n")

	return nil
}

func (g *gtool) AddGroup() error {
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
		log.Info(fmt.Sprintf("getting AD user by email: %s", email))
		query := fmt.Sprintf("mail=%s", email)
		user, err := g.searchEntry(query, []string{})
		if err != nil {
			return fmt.Errorf("failed to get user by email: %w", err)
		}

		users = append(users, user)
	}

	for _, groupCN := range groupCNs {
		log.Info(fmt.Sprintf("getting AD group by cn: %s", groupCN))
		query := fmt.Sprintf("cn=%s", groupCN)
		group, err := g.searchEntry(query, []string{})
		if err != nil {
			return fmt.Errorf("failed to get group by cn: %w", err)
		}

		groups = append(groups, group)
	}

	for _, user := range users {
		for _, group := range groups {
			log.Info(fmt.Sprintf("adding user %s to group %s", user.GetAttributeValue("mail"), group.GetAttributeValue("cn")))
			err := g.addUserToGroup(user, group)
			if err != nil {
				log.Error(fmt.Errorf("failed to add user to group: %w", err))
			}
		}
	}

	return nil
}
