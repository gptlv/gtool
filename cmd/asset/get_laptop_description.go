package asset

import (
	"log"

	"github.com/spf13/cobra"
)

var getLaptopDescriptionCmd = &cobra.Command{
	Use:   "get-laptop-description",
	Short: "Get description for a user's laptop",
	Run: func(cmd *cobra.Command, args []string) {
		err := assetHandler.PrintLaptopDescription()
		if err != nil {
			log.Fatal(err)
		}
	},
}
