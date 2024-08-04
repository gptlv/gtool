package issue

import (
	"fmt"
	"main/internal/interfaces"
	"main/jira-tools/cmd"

	"github.com/spf13/cobra"
)

var issueHandler interfaces.IssueHandler

func init() {
	injectIssueHandler()

	cmd.RootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(assignAllCmd)
	// issueCmd.AddCommand(deactivateInsightCmd)
	// issueCmd.AddCommand(grantPermissionCmd)
	// issueCmd.AddCommand(updateBlockTraineeCmd)
	// issueCmd.AddCommand(showEmptyComponentCmd)
	// issueCmd.AddCommand(addUserToGroupCmd)
	// issueCmd.AddCommand(dismissalOrHiringCmd)
	// issueCmd.AddCommand(disableActiveDirectoryCmd)
	// issueCmd.AddCommand(returnCCEquipmentCmd)
}

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'myapp issue --help' for more information on managing issues")
	},
}
