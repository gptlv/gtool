package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var grantAccessCmd = &cobra.Command{
	Use:   "grant-access",
	Short: "Process grant access issue",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.ProcessGrantAccessIssue()
		if err != nil {
			log.Fatal(err)
		}
	},
}
