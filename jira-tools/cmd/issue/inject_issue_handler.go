package issue

import (
	"main/internal/handlers"
	"main/internal/services"
	"main/jira-tools/cmd"
)

func injectIssueHandler() {
	issueService := services.NewIssueService(cmd.Client)
	activeDirectoryService := services.NewActiveDirectoryService(cmd.Conn)
	assetService := services.NewAssetService(cmd.Client)

	issueHandler = handlers.NewIssueHandler(issueService, activeDirectoryService, assetService)
}
