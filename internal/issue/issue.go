package issue

import "github.com/andygrunwald/go-jira"

type IssueService interface {
	GetAll(jql string) ([]jira.Issue, error)
	GetParent(issue *jira.Issue) (*jira.Issue, error)
	GetByID(ID string) (*jira.Issue, error)
	GetSubtaskByComponent(issue *jira.Issue, componentName string) (*jira.Subtasks, error)
	GetUserEmail(issue *jira.Issue) (string, error)
	Close(issue *jira.Issue) (*jira.Issue, error)
	BlockByIssue(currentIssue *jira.Issue, blockingIssue *jira.Issue) (*jira.Issue, error)
	PrintIssue(issue *jira.Issue)
	AssignIssueToMe(issue *jira.Issue) (*jira.Issue, error)
	WriteInternalComment(issue *jira.Issue, commentText string) (*jira.Comment, error)
}

const EMAIL_FIELD_KEY = "customfield_10145"

var internalCommentPayloadBody = `{
	"body": "%s",
	"properties": [
	  {
		"key": "sd.public.comment",
		"value": {
		   "internal": true
		}
	  }
	]
 }`

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

type BlockByIssuePayload struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
	Update struct {
		Issuelinks []struct {
			Add struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
				InwardIssue struct {
					Key string `json:"key"`
				} `json:"inwardIssue"`
			} `json:"add"`
		} `json:"issuelinks"`
	} `json:"update"`
}

type InternalCommentPayload struct {
	Body       string `json:"body"`
	Properties []struct {
		Key   string `json:"key"`
		Value struct {
			Internal bool `json:"internal"`
		} `json:"value"`
	} `json:"properties"`
}
