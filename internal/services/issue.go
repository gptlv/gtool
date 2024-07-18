package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/config"
	"main/internal/interfaces"
	"main/internal/models"

	"github.com/andygrunwald/go-jira"
)

type issueService struct {
	client *jira.Client
}

func NewIssueService(client *jira.Client) interfaces.IssueService {
	return &issueService{client: client}
}

func (issueService *issueService) GetAll(jql string) ([]jira.Issue, error) {
	fmt.Printf("Usecase: Running a JQL query '%s'\n", jql)
	issues, _, err := issueService.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 100}) //SearchOptions <- nil
	if err != nil {
		return []jira.Issue{}, fmt.Errorf("falied to get all issues: %w", err)
	}

	return issues, nil
}

func (issueService *issueService) Update(issue *jira.Issue, data map[string]interface{}) (*jira.Issue, error) {
	_, err := issueService.client.Issue.UpdateIssue(issue.ID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue %v: %w", issue.Key, err)
	}

	updatedIssue, err := issueService.GetByID(issue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated issue %v: %w", updatedIssue.Key, err)
	}

	return updatedIssue, nil

}

func (issueService *issueService) GetParent(issue *jira.Issue) (*jira.Issue, error) {
	parent := issue.Fields.Parent
	parentIssue, _, err := issueService.client.Issue.Get(parent.ID, nil)
	if err != nil {
		return parentIssue, err
	}

	return parentIssue, nil
}

func (issueService *issueService) GetByID(ID string) (*jira.Issue, error) {
	if ID == "" {
		return &jira.Issue{}, errors.New("empty ID")
	}

	issue, _, err := issueService.client.Issue.Get(ID, nil)
	if err != nil {
		return &jira.Issue{}, err
	}

	return issue, nil
}

func (issueService *issueService) GetSubtaskByComponent(issue *jira.Issue, componentName string) (*jira.Subtasks, error) {
	if componentName == "" {
		return nil, errors.New("empty component")
	}

	if issue == nil {
		return nil, errors.New("invalid issue")
	}

	subtasks := issue.Fields.Subtasks

	for _, st := range subtasks {
		issue, err := issueService.GetByID(st.ID)
		if err != nil {
			return nil, err
		}

		currentComponent := issue.Fields.Components[0].Name //unreliable

		if currentComponent == componentName {
			return st, nil
		}
	}

	return nil, errors.New("no such subtask")

}

// func (issueService *issueService) GetUserEmail(issue *jira.Issue) (string, error) {
// 	email, _ := issue.Fields.Unknowns.Value(config.EMAIL_FIELD_KEY)
// 	email = strings.TrimSpace(email.(string))

// 	return fmt.Sprintf("%v", email), nil
// }

func (issueService *issueService) Close(issue *jira.Issue) (*jira.Issue, error) {
	if issue == nil {
		return nil, errors.New("empty issue")
	}

	possibleTransitions, _, err := issueService.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "In Progress" {
			_, err = issueService.client.Issue.DoTransition(issue.ID, t.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	issue, err = issueService.GetByID(issue.ID)
	if err != nil {
		return nil, err
	}

	possibleTransitions, _, err = issueService.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Done" {
			_, err = issueService.client.Issue.DoTransition(issue.ID, t.ID)
			if err != nil {
				return nil, err
			}
		}

	}

	currentIssue, err := issueService.GetByID(issue.ID)
	if err != nil {
		return nil, err
	}

	return currentIssue, nil

}

func (issueService *issueService) BlockByIssue(currentIssue *jira.Issue, blockingIssue *jira.Issue) (*jira.Issue, error) {
	possibleTransitions, _, err := issueService.client.Issue.GetTransitions(currentIssue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Block" {
			body := fmt.Sprintf(config.BlockByIssuePayloadBody, t.ID, blockingIssue.Key)

			payload := new(models.BlockByIssuePayload)

			err := json.Unmarshal([]byte(body), payload)
			if err != nil {
				return nil, err
			}

			_, err = issueService.client.Issue.DoTransitionWithPayload(currentIssue.ID, &payload)
			if err != nil {
				return nil, err
			}
		}
	}

	updatedIssue, err := issueService.GetByID(currentIssue.ID)
	if err != nil {
		return nil, err
	}

	return updatedIssue, nil
}

func (issueService *issueService) AssignIssueToMe(issue *jira.Issue) (*jira.Issue, error) {
	users, _, err := issueService.client.User.Find("", jira.WithUsername("potlov.ga"))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	me := &users[0]

	_, err = issueService.client.Issue.UpdateAssignee(issue.ID, me)
	if err != nil {
		return nil, fmt.Errorf("failed to update assignee: %w", err)
	}

	updatedIssue, err := issueService.GetByID(issue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return updatedIssue, nil

}

func (issueService *issueService) PrintIssue(issue *jira.Issue) {
	fmt.Printf("%s (%s/%s): %+v\n", issue.Key, issue.Fields.Type.Name, issue.Fields.Priority.Name, issue.Fields.Summary)
}

func (issueService *issueService) WriteInternalComment(issue *jira.Issue, commentText string) (*jira.Comment, error) {
	body := fmt.Sprintf(config.InternalCommentPayloadBody, commentText)
	payload := new(models.InternalCommentPayload)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	endpoint := fmt.Sprintf("rest/api/2/issue/%v/comment", issue.ID)

	req, err := issueService.client.NewRequest("POST", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	comment := new(jira.Comment)
	resp, err := issueService.client.Do(req, comment)
	if err != nil {
		fmt.Printf("%+v", resp.Response)
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return comment, nil
}

func (issueService *issueService) Summarize(issue *jira.Issue) string {
	return fmt.Sprintf("[%s] %s", issue.Key, issue.Fields.Summary)
}

func (issueService *issueService) BlockUntilTomorrow(issue *jira.Issue) (*jira.Issue, error) {
	if issue == nil {
		return nil, errors.New("empty issue")
	}

	possibleTransitions, _, err := issueService.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Block" {
			// body := fmt.Sprintf(blockByIssuePayloadBody, t.ID, blockingIssue.Key)

			payload := new(models.BlockUntilTomorrowPayload)

			err := json.Unmarshal([]byte(config.BlockUntilTomorrowPayloadBody), payload)
			if err != nil {
				return nil, err
			}

			_, err = issueService.client.Issue.DoTransitionWithPayload(issue.ID, &payload)
			if err != nil {
				return nil, err
			}
		}
	}

	updatedIssue, err := issueService.GetByID(issue.ID)
	if err != nil {
		return nil, err
	}

	return updatedIssue, nil
}

func (issueService *issueService) GetUnresolvedSubtask(issue *jira.Issue) (*jira.Issue, error) {
	if issue == nil {
		return nil, errors.New("empty issue")
	}

	for _, subtask := range issue.Fields.Subtasks {
		substaskIssue, err := issueService.GetByID(subtask.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subtask issue %v by id: %w", substaskIssue.Key, err)
		}

		if substaskIssue.Fields.Status.ID != "6" {
			return substaskIssue, nil
		}
	}

	return nil, nil
}
