package asset

import (
	"fmt"

	"github.com/gptlv/gtools/common"
	"github.com/gptlv/gtools/internal/handlers"
	"github.com/gptlv/gtools/internal/interfaces"
	"github.com/gptlv/gtools/internal/services"
	"github.com/spf13/cobra"
)

func init() {
	initAssetHandler()
	AssetCmd.AddCommand(getLaptopDescriptionCmd)
	AssetCmd.AddCommand(generateRecordsCmd)
}

var assetHandler interfaces.AssetHandler

func initAssetHandler() error {
	client, err := common.GetJiraClient()
	if err != nil {
		return fmt.Errorf("failed to get jira client: %w", err)
	}

	assetService := services.NewAssetService(client)
	assetHandler = handlers.NewAssetHandler(assetService)

	return nil
}

var AssetCmd = &cobra.Command{
	Use:   "asset",
	Short: "Modify active directory entries",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'myapp issue --help' for more information on managing assets")
	},
}
