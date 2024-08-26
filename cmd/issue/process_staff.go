package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var processStaffCmd = &cobra.Command{
	Use:   "process-staff",
	Short: "Process staff issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.ProcessStaffIssues()
		if err != nil {
			log.Fatal(err)
		}
	},
}
