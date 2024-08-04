package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var assignAllCmd = &cobra.Command{
	Use:   "assign-all",
	Short: "Assign all deactivation issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.AssignAllDeactivateInsightIssuesToMe()
		if err != nil {
			log.Fatal(err)
		}
	},
}
