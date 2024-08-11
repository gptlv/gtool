package writeoff

import (
	"fmt"

	"github.com/gptlv/gtools/common"
	"github.com/gptlv/gtools/internal/handlers"
	"github.com/gptlv/gtools/internal/interfaces"
	"github.com/gptlv/gtools/internal/services"
	"github.com/spf13/cobra"
)

func init() {
	initWriteOffHandler()
	WriteOffCmd.AddCommand(generateRecordsCmd)

}

var writeOffHandler interfaces.WriteOffHandler

var WriteOffCmd = &cobra.Command{
	Use:   "write-off",
	Short: "Actions that have to do with writing off equipment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'myapp write-off --help' for more information on write-off actions")
	},
}

func initWriteOffHandler() error {
	client, err := common.GetJiraClient()
	if err != nil {
		return fmt.Errorf("failed to get jira client: %w", err)
	}

	writeOffService := services.NewWriteOffService()
	assetService := services.NewAssetService(client)

	writeOffHandler = handlers.NewWriteOffHandler(writeOffService, assetService)

	return nil
}
