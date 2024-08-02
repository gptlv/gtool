package issue

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DeactivateInsightCmd = &cobra.Command{
	Use:   "deactivate-insight",
	Short: "Process deactivate insight account issue",
	Run: func(cmd *cobra.Command, args []string) {
		issueHandler := getIssueHandler(cmd)

		err := issueHandler.ProcessDeactivateInsightAccountIssue()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}
