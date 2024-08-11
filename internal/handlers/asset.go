package handlers

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gptlv/gtools/internal/interfaces"
)

type assetHandler struct {
	assetService interfaces.AssetService
}

func NewAssetHandler(assetService interfaces.AssetService) interfaces.AssetHandler {
	return &assetHandler{assetService: assetService}
}

func (assetHandler *assetHandler) PrintLaptopDescription() error {
	var email string

	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	user, err := assetHandler.assetService.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	res, err := assetHandler.assetService.GetUserLaptops(user)
	if err != nil {
		log.Fatal(err)
	}

	for _, laptop := range res.ObjectEntries {
		description, err := assetHandler.assetService.GetLaptopDescription(&laptop)
		if err != nil {
			return fmt.Errorf("failed to get %v laptop description: %w", laptop.ObjectKey, err)
		}

		fmt.Printf("\n Ноутбук %s\n %s \n %s \n\n %s \n\n", description.Name, description.ISC, description.Serial, description.Cost)
	}

	return nil
}
