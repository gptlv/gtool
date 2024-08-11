package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var deactivateInsightCmd = &cobra.Command{
	Use:   "deactivate-insight",
	Short: "Process and deactivate Insight account issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.ProcessDeactivateInsightAccountIssues()
		if err != nil {
			log.Fatal(err)
		}
	},
}
