package controllers

import (
	"fmt"
	"log"
	"main/internal/interfaces"
)

type AssetHandler struct {
	assetService interfaces.AssetService
}

func NewAssetHandler(assetService interfaces.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (assetHandler *AssetHandler) PrintLaptopDescription() error {
	var email string

	fmt.Print("enter user's email: ")
	fmt.Scanln(&email)

	user, err := h.objectService.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	res, err := h.objectService.GetUserLaptops(user)
	if err != nil {
		log.Fatal(err)
	}

	for _, laptop := range res.ObjectEntries {
		description, err := h.objectService.GetLaptopDescription(&laptop)
		if err != nil {
			return fmt.Errorf("failed to get %v laptop description: %w", laptop.ObjectKey, err)
		}

		fmt.Printf("\n Ноутбук %s\n %s \n %s \n\n %s \n\n", description.Name, description.ISC, description.Serial, description.Cost)
	}

	return nil
}
