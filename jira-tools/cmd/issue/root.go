package issue

import (
	"fmt"
	"main/common"
	"main/internal/handlers"
	"main/internal/interfaces"
	"main/internal/services"

	"github.com/spf13/cobra"
)

var issueHandler interfaces.IssueHandler

func init() {
	initHandler()
	IssueCmd.AddCommand(assignAllCmd)
	IssueCmd.AddCommand(deactivateInsightCmd)
	IssueCmd.AddCommand(disableActiveDirectoryCmd)
	IssueCmd.AddCommand(grantAccessCmd)
	IssueCmd.AddCommand(processStaffCmd)
	IssueCmd.AddCommand(returnEquipmentCmd)
	IssueCmd.AddCommand(showEmptyCmd)
	IssueCmd.AddCommand(updateBlockTraineeCmd)
}

var IssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'myapp issue --help' for more information on managing issues")
	},
}

func initHandler() error {
	client, err := common.GetJiraClient()
	if err != nil {
		return fmt.Errorf("failed to get jira client: %w", err)
	}

	conn, err := common.GetLDAPConnection()
	if err != nil {
		return fmt.Errorf("failed to get ldap connection: %w", err)
	}

	issueService := services.NewIssueService(client)
	activeDirectoryService := services.NewActiveDirectoryService(conn)
	assetService := services.NewAssetService(client)

	issueHandler = handlers.NewIssueHandler(issueService, activeDirectoryService, assetService)

	return nil
}
