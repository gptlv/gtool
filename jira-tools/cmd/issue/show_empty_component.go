package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var ShowEmptyComponentCmd = &cobra.Command{
	Use:   "show-empty-component",
	Short: "Show all issues with empty component",
	Run: func(cmd *cobra.Command, args []string) {
		issueHandler := getIssueHandler(cmd)

		err := issueHandler.ShowIssuesWithEmptyComponent()
		if err != nil {
			log.Fatal(err)
		}
	},
}
