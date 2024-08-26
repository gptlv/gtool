package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var showEmptyCmd = &cobra.Command{
	Use:   "show-empty",
	Short: "Show issues with empty component",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.ShowIssuesWithEmptyComponent()
		if err != nil {
			log.Fatal(err)
		}
	},
}
