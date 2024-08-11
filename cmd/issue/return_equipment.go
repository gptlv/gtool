package issue

import (
	"log"

	"github.com/spf13/cobra"
)

var returnEquipmentCmd = &cobra.Command{
	Use:   "return-equipment",
	Short: "Process return CC equipment issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := issueHandler.ProcessReturnCCEquipmentIssues()
		if err != nil {
			log.Fatal(err)
		}
	},
}
