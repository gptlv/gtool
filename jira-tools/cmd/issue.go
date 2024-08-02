package cmd

import (
	"fmt"
	"main/jira-tools/cmd/issue"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issue.DeactivateInsightCmd)
	issueCmd.AddCommand(issue.AssignAllCmd)
	issueCmd.AddCommand(issue.GrantPermissionCmd)
	issueCmd.AddCommand(issue.UpdateBlockTraineeCmd)
	issueCmd.AddCommand(issue.ShowEmptyComponentCmd)
	issueCmd.AddCommand(issue.AddUserToGroupCmd)
	issueCmd.AddCommand(issue.DismissalOrHiringCmd)
	issueCmd.AddCommand(issue.DisableActiveDirectoryCmd)
	issueCmd.AddCommand(issue.ReturnCCEquipmentCmd)
}

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'myapp issue --help' for more information on managing issues")
	},
}
