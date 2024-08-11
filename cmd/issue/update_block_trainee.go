package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var updateBlockTraineeCmd = &cobra.Command{
	Use:   "update-block-trainee",
	Short: "Update block trainee issue",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.UpdateBlockTraineeIssue()
		if err != nil {
			log.Fatal(err)
		}
	},
}
