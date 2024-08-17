package handlers

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gocarina/gocsv"
	"github.com/gptlv/gtools/internal/domain"
	"github.com/gptlv/gtools/internal/interfaces"
	"github.com/gptlv/gtools/util"
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

func (assetHandler *assetHandler) GenerateWriteOffRecords() error {
	assets := []*domain.Asset{}

	inputFile := "info.csv"
	outputFile := "write_off_records.csv"

	log.Info(fmt.Sprintf("input file name: %v", inputFile))
	log.Info(fmt.Sprintf("output file name: %v\n", outputFile))

	log.Info(fmt.Sprintf("reading csv input file %v\n", inputFile))
	input, err := util.ReadInputFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}
	defer input.Close()

	log.Info("unmarshaling input file\n")
	err = gocsv.UnmarshalFile(input, &assets)
	if err != nil {
		return fmt.Errorf("failed to unmarshal input file %s: %w", inputFile, err)
	}

	for _, asset := range assets {
		log.SetPrefix(asset.ISC)
		log.Info("getting laptop by isc")
		laptop, err := assetHandler.assetService.GetByISC(asset.ISC)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		log.Info("getting laptop description")
		description, err := assetHandler.assetService.GetLaptopDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		asset.Name = description.Name
		asset.InventoryID = description.InventoryID
		asset.Serial = description.Serial

		log.Info("finished processing the asset\n")
		log.SetPrefix("")
	}

	log.Info("opening output file")
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	log.Info("marshaling output file")
	err = gocsv.MarshalFile(&assets, output)
	if err != nil {
		return fmt.Errorf("failed to unmarshal csv file: %w", err)
	}

	log.Info("finished generating write-off records\n")

	return nil
}
