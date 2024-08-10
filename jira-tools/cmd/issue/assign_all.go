package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var assignAllCmd = &cobra.Command{
	Use:   "assign-all",
	Short: "Assign all automatable issues to current user",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.AssignAutomatableIssuesToCurrentUser()
		if err != nil {
			log.Fatal(err)
		}
	},
}
