package issue

import (
	"fmt"

	"github.com/gptlv/gtools/common"
	"github.com/gptlv/gtools/internal/handlers"
	"github.com/gptlv/gtools/internal/interfaces"
	"github.com/gptlv/gtools/internal/services"
	"github.com/spf13/cobra"
)

var issueHandler interfaces.IssueHandler

func init() {
	initIssueHandler()
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

func initIssueHandler() error {
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
