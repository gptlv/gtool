package writeoff

import (
	"log"

	"github.com/spf13/cobra"
)

var generateRecordsCmd = &cobra.Command{
	Use:   "generate-records",
	Short: "generate write-off records",
	Run: func(cmd *cobra.Command, args []string) {
		err := writeOffHandler.GenerateWriteOffRecords()
		if err != nil {
			log.Fatal(err)
		}
	},
}
