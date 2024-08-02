package issue

import (
	"main/internal/handlers"
	"main/internal/interfaces"
	"main/internal/services"

	"github.com/andygrunwald/go-jira"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
)

func getIssueHandler(cmd *cobra.Command) interfaces.IssueHandler {
	ctx := cmd.Context()
	client := ctx.Value("client").(*jira.Client)
	conn := ctx.Value("conn").(*ldap.Conn)

	issueService := services.NewIssueService(client)
	assetService := services.NewAssetService(client)
	activeDirectoryService := services.NewActiveDirectoryService(conn)

	issueHandler := handlers.NewIssueHandler(issueService, activeDirectoryService, assetService)

	return issueHandler
}
