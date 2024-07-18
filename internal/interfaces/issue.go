package interfaces

import "github.com/andygrunwald/go-jira"

type IssueService interface {
	GetAll(jql string) ([]jira.Issue, error)
	GetParent(issue *jira.Issue) (*jira.Issue, error)
	GetByID(ID string) (*jira.Issue, error)
	GetSubtaskByComponent(issue *jira.Issue, componentName string) (*jira.Subtasks, error)
	GetUserEmail(issue *jira.Issue) (string, error)
	Update(issue *jira.Issue, data map[string]interface{}) (*jira.Issue, error)
	Close(issue *jira.Issue) (*jira.Issue, error)
	BlockByIssue(currentIssue *jira.Issue, blockingIssue *jira.Issue) (*jira.Issue, error)
	PrintIssue(issue *jira.Issue)
	AssignIssueToMe(issue *jira.Issue) (*jira.Issue, error)
	WriteInternalComment(issue *jira.Issue, commentText string) (*jira.Comment, error)
	Summarize(issue *jira.Issue) string
	BlockUntilTomorrow(issue *jira.Issue) (*jira.Issue, error)
	GetUnresolvedSubtask(issue *jira.Issue) (*jira.Issue, error)
}

type IssueHandler interface {
	DeactivateInsight() error
	GrantPermission() error
	UpdateBlockTraineeIssue() error
	ShowIssuesWithEmptyComponent() error
	BlockDismissedUserInActiveDirectory() error
}
