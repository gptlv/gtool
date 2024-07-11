package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
)

type issueService struct {
	client *jira.Client
}

func NewIssueService(client *jira.Client) IssueService {
	return &issueService{client: client}
}

func (s *issueService) GetAll(jql string) ([]jira.Issue, error) {
	fmt.Printf("Usecase: Running a JQL query '%s'\n", jql)
	issues, _, err := s.client.Issue.Search(jql, &jira.SearchOptions{MaxResults: 100}) //SearchOptions <- nil
	if err != nil {
		return []jira.Issue{}, fmt.Errorf("falied to get all issues: %w", err)
	}

	return issues, nil
}

func (s *issueService) Update(issue *jira.Issue, data map[string]interface{}) (*jira.Issue, error) {
	_, err := s.client.Issue.UpdateIssue(issue.ID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue %v: %w", issue.Key, err)
	}

	updatedIssue, err := s.GetByID(issue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated issue %v: %w", updatedIssue.Key, err)
	}

	return updatedIssue, nil

}

func (s *issueService) GetParent(issue *jira.Issue) (*jira.Issue, error) {
	parent := issue.Fields.Parent
	parentIssue, _, err := s.client.Issue.Get(parent.ID, nil)
	if err != nil {
		return parentIssue, err
	}

	return parentIssue, nil
}

func (s *issueService) GetByID(ID string) (*jira.Issue, error) {
	if ID == "" {
		return &jira.Issue{}, errors.New("empty ID")
	}

	issue, _, err := s.client.Issue.Get(ID, nil)
	if err != nil {
		return &jira.Issue{}, err
	}

	return issue, nil
}

func (s *issueService) GetSubtaskByComponent(issue *jira.Issue, componentName string) (*jira.Subtasks, error) {
	if componentName == "" {
		return nil, errors.New("empty component")
	}

	if issue == nil {
		return nil, errors.New("invalid issue")
	}

	subtasks := issue.Fields.Subtasks

	for _, st := range subtasks {
		issue, err := s.GetByID(st.ID)
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

func (s *issueService) GetUserEmail(issue *jira.Issue) (string, error) {
	email, _ := issue.Fields.Unknowns.Value(EMAIL_FIELD_KEY)
	email = strings.TrimSpace(email.(string))

	return fmt.Sprintf("%v", email), nil
}

func (s *issueService) Close(issue *jira.Issue) (*jira.Issue, error) {
	if issue == nil {
		return nil, errors.New("empty issue")
	}

	possibleTransitions, _, err := s.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "In Progress" {
			fmt.Printf("[%v] %v -> %v\n", issue.Key, issue.Fields.Status.Name, t.Name)
			_, err = s.client.Issue.DoTransition(issue.ID, t.ID)
		}

		if err != nil {
			return nil, err
		}
	}

	issue, err = s.GetByID(issue.ID)
	if err != nil {
		return nil, err
	}

	possibleTransitions, _, err = s.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Done" {
			fmt.Printf("[%v] %v -> %v\n", issue.Key, issue.Fields.Status.Name, t.Name)
			_, err = s.client.Issue.DoTransition(issue.ID, t.ID)
			if err != nil {
				return nil, err
			}
		}

	}

	currentIssue, err := s.GetByID(issue.ID)
	if err != nil {
		return nil, err
	}

	return currentIssue, nil

}

func (s *issueService) BlockByIssue(currentIssue *jira.Issue, blockingIssue *jira.Issue) (*jira.Issue, error) {
	possibleTransitions, _, err := s.client.Issue.GetTransitions(currentIssue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Block" {
			body := fmt.Sprintf(blockByIssuePayloadBody, t.ID, blockingIssue.Key)

			payload := new(BlockByIssuePayload)

			err := json.Unmarshal([]byte(body), payload)
			if err != nil {
				return nil, err
			}

			_, err = s.client.Issue.DoTransitionWithPayload(currentIssue.ID, &payload)
			if err != nil {
				return nil, err
			}
		}
	}

	updatedIssue, err := s.GetByID(currentIssue.ID)
	if err != nil {
		return nil, err
	}

	return updatedIssue, nil
}

func (s *issueService) AssignIssueToMe(issue *jira.Issue) (*jira.Issue, error) {
	users, _, err := s.client.User.Find("", jira.WithUsername("potlov.ga"))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	me := &users[0]

	_, err = s.client.Issue.UpdateAssignee(issue.ID, me)
	if err != nil {
		return nil, fmt.Errorf("failed to update assignee: %w", err)
	}

	updatedIssue, err := s.GetByID(issue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return updatedIssue, nil

}

func (s *issueService) PrintIssue(issue *jira.Issue) {
	fmt.Printf("%s (%s/%s): %+v\n", issue.Key, issue.Fields.Type.Name, issue.Fields.Priority.Name, issue.Fields.Summary)
}

func (s *issueService) WriteInternalComment(issue *jira.Issue, commentText string) (*jira.Comment, error) {
	body := fmt.Sprintf(internalCommentPayloadBody, commentText)
	payload := new(InternalCommentPayload)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	endpoint := fmt.Sprintf("rest/api/2/issue/%v/comment", issue.ID)

	req, err := s.client.NewRequest("POST", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	comment := new(jira.Comment)
	resp, err := s.client.Do(req, comment)
	if err != nil {
		fmt.Printf("%+v", resp.Response)
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return comment, nil
}

func (s *issueService) Summarize(issue *jira.Issue) string {
	return fmt.Sprintf("[%s] %s", issue.Key, issue.Fields.Summary)
}

func (s *issueService) BlockUntilTomorrow(issue *jira.Issue) (*jira.Issue, error) {
	if issue == nil {
		return nil, errors.New("empty issue")
	}

	possibleTransitions, _, err := s.client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Block" {
			// body := fmt.Sprintf(blockByIssuePayloadBody, t.ID, blockingIssue.Key)

			payload := new(BlockUntilTomorrowPayload)

			err := json.Unmarshal([]byte(BlockUntilTomorrowPayloadBody), payload)
			if err != nil {
				return nil, err
			}

			_, err = s.client.Issue.DoTransitionWithPayload(issue.ID, &payload)
			if err != nil {
				return nil, err
			}
		}
	}

	updatedIssue, err := s.GetByID(issue.ID)
	if err != nil {
		return nil, err
	}

	return updatedIssue, nil
}
