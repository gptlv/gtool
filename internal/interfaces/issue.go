package interfaces

import (
	"main/internal/domain"
)

type IssueService interface {
	GetAll(jql string) ([]*domain.Issue, error)
	GetParent(issue *domain.Issue) (*domain.Issue, error)
	GetByID(ID string) (*domain.Issue, error)
	GetSubtaskByComponent(issue *domain.Issue, componentName string) (*domain.Subtasks, error)
	GetUserEmail(issue *domain.Issue) (string, error)
	Update(issue *domain.Issue, data map[string]interface{}) (*domain.Issue, error)
	Close(issue *domain.Issue) (*domain.Issue, error)
	BlockByIssue(currentIssue *domain.Issue, blockingIssue *domain.Issue) (*domain.Issue, error)
	PrintIssue(issue *domain.Issue)
	AssignIssueToMe(issue *domain.Issue) (*domain.Issue, error)
	// WriteInternalComment(issue *domain.Issue, commentText string) (*domain.Comment, error)
	Summarize(issue *domain.Issue) string
	BlockUntilTomorrow(issue *domain.Issue) (*domain.Issue, error)
}

type IssueUsecase interface {
	DeactivateInsight() error
	GrantPermission() error
	UpdateBlockTraineeIssue() error
	ShowIssuesWithEmptyComponent() error
	BlockDismissedUserInActiveDirectory() error
}
