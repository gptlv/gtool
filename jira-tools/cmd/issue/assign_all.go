package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var AssignAllCmd = &cobra.Command{
	Use:   "assign-all",
	Short: "Assign all deactivation issues",
	Run: func(cmd *cobra.Command, args []string) {
		issueHandler := getIssueHandler(cmd)

		err := issueHandler.AssignAllDeactivateInsightIssuesToMe()
		if err != nil {
			log.Fatal(err)
		}

	},
}
