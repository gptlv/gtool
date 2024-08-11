package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var disableActiveDirectoryCmd = &cobra.Command{
	Use:   "disable-ad",
	Short: "Process disable active directory account issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.ProcessDisableActiveDirectoryAccountIssues()
		if err != nil {
			log.Fatal(err)
		}
	},
}
