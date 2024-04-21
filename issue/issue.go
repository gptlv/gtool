package issue

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"

	"github.com/andygrunwald/go-jira"
)

const EMAIL_FIELD_KEY = "customfield_10145"

var blockByIssuePayloadBody = `
{
	"transition": {
		"id": "%v"
	},
	"update": {
		"issuelinks": [
			{
				"add": {
					"type": {
						"name": "Blocks"
					},
					"inwardIssue": {
						"key": "%v"
					}
				}
			}
		]
	}
}`

func GetAll(client *jira.Client, jql string) ([]jira.Issue, error) {
	fmt.Printf("Usecase: Running a JQL query '%s'\n", jql)
	issues, _, err := client.Issue.Search(jql, nil)
	if err != nil {
		return []jira.Issue{}, fmt.Errorf("falied to get all issues: %w", err)
	}

	return issues, nil
}

func GetParent(client *jira.Client, issue *jira.Issue) (*jira.Issue, error) {
	parent := issue.Fields.Parent
	parentIssue, _, err := client.Issue.Get(parent.ID, nil)
	if err != nil {
		return parentIssue, err
	}

	return parentIssue, nil
}

func GetByID(client *jira.Client, ID string) (*jira.Issue, error) {
	if ID == "" {
		return &jira.Issue{}, errors.New("empty ID")
	}

	issue, _, err := client.Issue.Get(ID, nil)
	if err != nil {
		return &jira.Issue{}, err
	}

	return issue, nil
}

func GetSubtaskByComponent(client *jira.Client, issue *jira.Issue, componentName string) (*jira.Subtasks, error) {
	if componentName == "" {
		return nil, errors.New("empty component")
	}

	if issue == nil {
		return nil, errors.New("invalid issue")
	}

	subtasks := issue.Fields.Subtasks

	for _, st := range subtasks {
		issue, err := GetByID(client, st.ID)
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

func GetUserEmail(client *jira.Client, key string) (string, error) {
	issueEndPoint := fmt.Sprintf("rest/api/2/issue/%s", key)
	req, _ := client.NewRequest("GET", issueEndPoint, nil)

	issue := new(jira.Issue)

	_, err := client.Do(req, issue)
	if err != nil {
		return "", err
	}

	email, _ := issue.Fields.Unknowns.Value(EMAIL_FIELD_KEY)

	return fmt.Sprintf("%v", email), nil
}

func Close(client *jira.Client, issue *jira.Issue) (*jira.Issue, error) {
	if issue == nil {
		return nil, errors.New("empty issue")
	}

	possibleTransitions, _, err := client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "In Progress" {
			fmt.Printf("[%v] %v -> %v\n", issue.Key, issue.Fields.Status.Name, t.Name)
			_, err = client.Issue.DoTransition(issue.ID, t.ID)
		}

		if err != nil {
			return nil, err
		}
	}

	issue, err = GetByID(client, issue.ID)
	if err != nil {
		return nil, err
	}

	possibleTransitions, _, err = client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Done" {
			fmt.Printf("[%v] %v -> %v\n", issue.Key, issue.Fields.Status.Name, t.Name)
			_, err = client.Issue.DoTransition(issue.ID, t.ID)
			if err != nil {
				return nil, err
			}
		}

	}

	currentIssue, err := GetByID(client, issue.ID)
	if err != nil {
		return nil, err
	}

	return currentIssue, nil

}

func BlockByIssue(client *jira.Client, currentIssue *jira.Issue, blockingIssue *jira.Issue) (*jira.Issue, error) {
	possibleTransitions, _, err := client.Issue.GetTransitions(currentIssue.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range possibleTransitions {
		if t.Name == "Block" {
			body := fmt.Sprintf(blockByIssuePayloadBody, t.ID, blockingIssue.Key)

			payload := new(types.BlockByIssuePayload)

			err := json.Unmarshal([]byte(body), payload)
			if err != nil {
				return nil, err
			}

			_, err = client.Issue.DoTransitionWithPayload(currentIssue.ID, &payload)
			if err != nil {
				return nil, err
			}
		}
	}

	updatedIssue, err := GetByID(client, currentIssue.ID)
	if err != nil {
		return nil, err
	}

	return updatedIssue, nil
}

func OutputResponse(issues []jira.Issue, resp *jira.Response) {
	fmt.Printf("Call to %s\n", resp.Request.URL)
	fmt.Printf("Response Code: %d\n", resp.StatusCode)
	fmt.Println("==================================")
	for _, i := range issues {
		PrintIssue(&i)
	}
}

func PrintIssue(issue *jira.Issue) {
	fmt.Printf("%s (%s/%s): %+v\n", issue.Key, issue.Fields.Type.Name, issue.Fields.Priority.Name, issue.Fields.Summary)
}
