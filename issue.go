package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
)

var internalCommentProperty = jira.Property{Key: "sd.public.comment", Value: jira.Value{Internal: true}}

func (g *gtool) processDeactivateInsightIssue(issue *jira.Issue) error {
	log.SetPrefix(issue.Key)

	comment := new(jira.Comment)
	comment.Body = "[https://wiki.sbmt.io/x/sPjivQ]"
	comment.Properties = []jira.Property{internalCommentProperty}

	summaryFields := strings.Fields(issue.Fields.Summary)
	userEmail := summaryFields[len(summaryFields)-1]

	log.Info("getting parent issue")
	parentIssue, _, err := g.client.Issue.Get(issue.Fields.Parent.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to get parent issue %s: %w", issue.Fields.Parent.Key, err)
	}

	log.Info(fmt.Sprintf("found parent issue: %s", parentIssue.Fields.Summary))
	log.Info("getting return equipment subtask")

	component := new(jira.Component)
	component.Name = "Возврат оборудования"
	returnEquipmentIssue, err := g.getSubtaskByComponent(parentIssue, component)
	if err != nil {
		return fmt.Errorf("failed to get subtask by component: %w", err)
	}

	if returnEquipmentIssue.Fields.Status.Name != "Closed" {
		log.Info("return equipment subtask is unresolved")
		log.Info("blocking the issue by return equipment subtask")

		_, _, err := g.blockByIssue(issue, returnEquipmentIssue)
		if err != nil {
			return fmt.Errorf("failed to block issue by subtask")
		}

		return nil
	}

	log.Info("return equipment subtask is resolved")
	userList, _, err := g.getUserByEmail(userEmail)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(userList.ObjectEntries) == 0 {
		log.Info(fmt.Sprintf("found no user with email %s", userEmail))
		comment.Body = fmt.Sprintf("Нет пользователя c email %s в Insight", userEmail)
		return fmt.Errorf("no such user")
	}

	user := &userList.ObjectEntries[0]

	laptopList, _, err := g.getUserLaptops(user)
	if err != nil {
		return fmt.Errorf("failed to get user laptops: %w", err)
	}

	if len(laptopList.ObjectEntries) != 0 {
		log.Info("user still has attached laptops, skipping the issue")
		return nil
	}

	log.Info("user has no attached laptops")

	log.Info("setting user category to BYOD")
	_, _, err = g.setUserCategory(user, "BYOD")
	if err != nil {
		return fmt.Errorf("failed to set user category: %w", err)
	}

	log.Info("setting user status to inactive")
	_, _, err = g.disableUser(user)
	if err != nil {
		return fmt.Errorf("failed to disable user: %w", err)
	}

	log.Info("adding internal comment")
	_, _, err = g.client.Issue.AddComment(issue.ID, comment)
	if err != nil {
		return fmt.Errorf("failed to add internal comment: %w", err)
	}

	transitionChain := []string{"In Progress", "Done"}
	_, _, err = g.doTransitionChain(issue, transitionChain)
	if err != nil {
		return fmt.Errorf("failed to do transition chain: %w", err)
	}

	return nil
}

func (g *gtool) processStaffIssue(issue *jira.Issue) error {
	comment := new(jira.Comment)
	comment.Properties = []jira.Property{internalCommentProperty}
	for _, component := range issue.Fields.Components {
		if component.Name == "Dismissal" {
			comment.Body = "[https://wiki.sbmt.io/x/jgeLvg]"
		}

		if component.Name == "Hiring" {
			comment.Body = "[https://wiki.sbmt.io/x/ogeLvg]"
		}
	}

	log.SetPrefix(issue.Key)
	log.Info("looping over subtasks")
	unresolvedSubtask, err := g.getUnresolvedSubtask(issue)
	if err != nil {
		return fmt.Errorf("failed to get unresolved subtask: %w", err)
	}

	if unresolvedSubtask != nil {
		log.Info(fmt.Sprintf("found unresolved subtask %s", issue))
		log.Info("blocking the issue by unresolved subtask")

		_, _, err := g.blockByIssue(issue, unresolvedSubtask)
		if err != nil {
			return fmt.Errorf("failed to block by issue: %w", err)
		}

		return nil
	}

	log.Info("all subtasks are resolved")
	log.Info("adding internal comment")
	_, _, err = g.client.Issue.AddComment(issue.ID, comment)
	if err != nil {
		return fmt.Errorf("failed to add internal comment: %w", err)
	}

	transitionChain := []string{"In Progress", "Done"}
	_, _, err = g.doTransitionChain(issue, transitionChain)
	if err != nil {
		return fmt.Errorf("failed to do transition chain: %w", err)
	}

	return nil
}

func (g *gtool) updateSummary(issue *jira.Issue, trailingText string) (*jira.Issue, *jira.Response, error) {
	currentSummary := strings.TrimSpace(issue.Fields.Summary)
	newSummary := currentSummary + " " + trailingText

	payload := &jira.Issue{} //TODO test new(T) vs &T{}
	payload.Key = issue.Key
	payload.Fields = &jira.IssueFields{Summary: newSummary}

	return g.client.Issue.Update(payload)
}

func (g *gtool) blockByIssue(issue *jira.Issue, blockingIssue *jira.Issue) (*jira.Issue, *jira.Response, error) {
	transitions, _, err := g.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transitions: %w", err)
	}

	link := new(jira.IssueLink)
	link.Type.ID = "10000"
	link.OutwardIssue = &jira.Issue{ID: issue.ID}
	link.InwardIssue = &jira.Issue{ID: blockingIssue.ID}

	_, err = g.client.Issue.AddLink(link)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update issue: %w", err)
	}

	for _, transition := range transitions {
		if transition.Name == "Block" {
			_, err := g.client.Issue.DoTransition(issue.ID, transition.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to block the issue: %w", err)
			}
		}
	}

	return g.client.Issue.Get(issue.ID, nil)
}

func (g *gtool) getSubtaskByComponent(issue *jira.Issue, component *jira.Component) (*jira.Issue, error) {
	for _, subtask := range issue.Fields.Subtasks {
		subtaskIssue, _, err := g.client.Issue.Get(subtask.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get subtask issue: %w", err)
		}

		for _, currentComponent := range subtaskIssue.Fields.Components {
			if currentComponent.Name == component.Name {
				return subtaskIssue, nil
			}
		}
	}

	return nil, errors.New("no such subtask")
}

func (g *gtool) getUnresolvedSubtask(issue *jira.Issue) (*jira.Issue, error) {
	for _, subtask := range issue.Fields.Subtasks {
		if subtask.Fields.Status.Name != "Closed" {
			issue, _, err := g.client.Issue.Get(subtask.ID, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get issue: %w", err)
			}

			return issue, nil
		}
	}

	return nil, nil
}

func (g *gtool) doTransitionChain(issue *jira.Issue, transitions []string) (*jira.Issue, *jira.Response, error) {
	options := new(jira.GetQueryOptions)
	options.Expand = "transitions"

	for _, target := range transitions {
		issue, _, err := g.client.Issue.Get(issue.ID, options)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get issue: %w", err)
		}

		for _, issueTransition := range issue.Transitions {
			if issueTransition.Name != target {
				continue
			}

			log.Info(fmt.Sprintf("doing transition [%s] -> [%s]", issue.Fields.Status.Name, target))
			_, err = g.client.Issue.DoTransition(issue.ID, issueTransition.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to do transition")
			}
		}
	}

	return g.client.Issue.Get(issue.ID, options)
}

func (g *gtool) processReturnEquipmentIssue(issue *jira.Issue) error {
	noEquipmentjobTitles := []string{"оператор"}

	comment := new(jira.Comment)
	comment.Body = "Оборудование не отправлялось"
	comment.Properties = []jira.Property{internalCommentProperty}

	log.SetPrefix(issue.Key)

	summaryFields := strings.Fields(issue.Fields.Summary)
	email := summaryFields[len(summaryFields)-1]

	log.Info(fmt.Sprintf("getting user by email %v", email))
	userList, _, err := g.getUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user %v by email: %w", email, err)
	}

	if len(userList.ObjectEntries) == 0 {
		log.Info(fmt.Sprintf("no user with email %s", email))
		comment.Body = fmt.Sprintf("Нет пользователя c email %s в Insight", email)

		log.Info(fmt.Sprintf("adding internal comment: %s", comment.Body))
		_, _, err := g.client.Issue.AddComment(issue.ID, comment)
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}

		_, err = g.decline(issue)
		if err != nil {
			return fmt.Errorf("failed to decline the isssue: %w", err)
		}

		return nil
	}

	user := &userList.ObjectEntries[0]

	laptopList, _, err := g.getUserLaptops(user)
	if err != nil {
		return fmt.Errorf("failed to get user laptops")
	}

	if len(laptopList.ObjectEntries) != 0 {
		log.Info("user still has attached laptops")
		return nil
	}

	log.Info("getting parent issue")
	dismissalIssue, _, err := g.client.Issue.Get(issue.Fields.Parent.ID, nil)
	if err != nil {
		return fmt.Errorf("faield to get parent issue: %w", err)
	}

	title, err := dismissalIssue.Fields.Unknowns.String("customfield_10197")
	if err != nil {
		return fmt.Errorf("failed to get title: %w", err)
	}

	title = strings.TrimSpace(strings.ToLower(title))

	for _, noEquipmentTitle := range noEquipmentjobTitles {
		if title == noEquipmentTitle {
			_, err := g.decline(issue)
			if err != nil {
				return fmt.Errorf("failed to decline the isssue: %w", err)
			}
		}
	}

	log.Info("no matching titles")
	return nil
}

func (g *gtool) decline(issue *jira.Issue) (*jira.Response, error) {
	payload := new(jira.CreateTransitionPayload)
	transitions, _, err := g.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}
	for _, transition := range transitions {
		if transition.Name == "Closed" {
			payload.Transition.ID = transition.ID
			payload.Fields.Resolution.Name = "Won't Do"
		}
	}

	return g.client.Issue.DoTransitionWithPayload(issue.ID, payload)
}

func (g *gtool) processDisableActiveDirectoryIssue(issue *jira.Issue) error {
	log.SetPrefix(issue.Key)
	var comments []*jira.Comment

	blockComment := new(jira.Comment)
	blockComment.Body = "Заблокировал idm"
	blockComment.Properties = []jira.Property{internalCommentProperty}

	wikiComment := new(jira.Comment)
	wikiComment.Body = "[https://wiki.sbmt.io/x/mB6zsg]"
	blockComment.Properties = []jira.Property{internalCommentProperty}

	comments = append(comments, blockComment, wikiComment)

	summary := strings.Split(strings.TrimSpace(issue.Fields.Summary), " ")
	email := summary[len(summary)-1]

	log.Info(fmt.Sprintf("getting active directory user by email: %s", email))
	userQuery := fmt.Sprintf("mail=%s", email)
	user, err := g.searchEntry(userQuery, []string{"userAccountControl"})
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	value := user.GetAttributeValue("userAccountControl") //512 -- active, 514 -- inactive
	userAccountControlValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("failed to convert userAccountControl attribute to int: %w", err)
	}

	if userAccountControlValue != 514 {
		log.Info("user is active")
		log.Info("blocking the issue until tomorrow")
		tomorrow := time.Now().AddDate(0, 0, 1)
		_, err := g.blockUntil(issue, tomorrow)
		if err != nil {
			return fmt.Errorf("failed to block issue until tomorrow: %w", err)
		}

		return nil
	}

	log.Info("user is inactive")
	log.Info("closing the issue")

	for _, comment := range comments {
		log.Info(fmt.Sprintf("adding internal comment \"%s\"", comment.Body))
		_, _, err := g.client.Issue.AddComment(issue.ID, comment)
		if err != nil {
			return fmt.Errorf("failed to add comment \"%s\": %w", comment.Body, err)
		}
	}

	chain := []string{"In Progress", "Done"}
	_, _, err = g.doTransitionChain(issue, chain)
	if err != nil {
		return fmt.Errorf("failed to do transition chain: %w", err)
	}

	return nil
}

func (g *gtool) blockUntil(issue *jira.Issue, date time.Time) (*jira.Response, error) {
	var blockTransition jira.Transition
	payload := new(jira.CreateTransitionPayload)
	formattedDate := date.Format("2006-01-02")
	transitions, _, err := g.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}

	for _, transition := range transitions {
		if transition.Name == "Block" {
			blockTransition = transition //? legit
		}
	}

	payload.Transition = jira.TransitionPayload{ID: blockTransition.ID}
	payload.Fields.BlockUntil = formattedDate

	return g.client.Issue.DoTransitionWithPayload(issue.ID, payload)

}

func (g *gtool) processGrantAccessIssue(issue *jira.Issue) error {
	var groupCN string
	log.SetPrefix(issue.Key)

	comment := new(jira.Comment)
	comment.Body = "[https://wiki.sbmt.io/x/WcPivQ]"
	comment.Properties = []jira.Property{{Key: "sd.public.comment", Value: jira.Value{Internal: true}}}

	summary := strings.Fields(issue.Fields.Summary)
	email := summary[len(summary)-1]

	roleField := "customfield_13063"
	role, err := issue.Fields.Unknowns.StringSlice(roleField)
	if err != nil {
		return fmt.Errorf("failed to get role field: %w", err)
	}

	if len(role) < 1 {
		return fmt.Errorf("empty role field")
	}

	re := regexp.MustCompile(`\((IR-\d+)\)`)

	matches := re.FindStringSubmatch(role[0])
	if len(matches) < 2 {
		return fmt.Errorf("identifier not found in input string")
	}

	roleID := matches[1]

	log.Info(fmt.Sprintf("getting information resource by key %s", roleID))
	object, _, err := g.client.Object.Get(roleID, nil)
	if err != nil {
		return fmt.Errorf("failed to get role obejct: %w", err)
	}

	for _, attribute := range object.Attributes {
		if attribute.ObjectTypeAttributeID == 8527 { //? or use name instead
			groupCN = strings.TrimSpace(attribute.ObjectAttributeValues[0].Value)
		}
	}

	log.Info(fmt.Sprintf("getting group by CN %s", groupCN))
	groupQuery := fmt.Sprintf("cn=%s", groupCN)
	group, err := g.searchEntry(groupQuery, []string{})
	if err != nil {
		return fmt.Errorf("failed to get group by CN")
	}

	log.Info("getting user by email")
	userQuery := fmt.Sprintf("(mail=%s)", email)
	user, err := g.searchEntry(userQuery, []string{})
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	log.Info("adding user to group")
	err = g.addUserToGroup(user, group)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	log.Info("adding internal comment")
	_, _, err = g.client.Issue.AddComment(issue.ID, comment)
	if err != nil {
		return fmt.Errorf("failed to add internal comment: %w", err)
	}

	return nil
}

func (g *gtool) processBlockTraineeIssue(issue *jira.Issue) error {
	log.SetPrefix(issue.Key)

	description := issue.Fields.Description
	descriptionFields := strings.Fields(description)
	if len(descriptionFields) < 5 {
		return fmt.Errorf("incorrect description")
	}

	email := descriptionFields[4] //Заблокировать доступы стажеру оператору {email} в связи с тем, что он не завершил обучение

	log.Info(fmt.Sprintf("user email: %s", email))

	for _, subtask := range issue.Fields.Subtasks {
		issue, _, err := g.client.Issue.Get(subtask.ID, nil)
		if err != nil {
			return fmt.Errorf("failed to get subtask issue: %w", err)
		}

		log.Info(fmt.Sprintf("updating summary for %s", issue))
		_, _, err = g.updateSummary(issue, email)
		if err != nil {
			return fmt.Errorf("failed to update issue summary: %w", err)
		}
	}

	return nil
}
