package handlers

import (
	"fmt"
	"main/internal/interfaces"
	"main/internal/models"
)

type dismissalHandler struct {
	dismissalService interfaces.DismissalService
	assetService     interfaces.AssetService
}

func NewdismissalHandler(dismissalService interfaces.DismissalService, assetService interfaces.AssetService) *dismissalHandler {
	return &dismissalHandler{dismissalService: dismissalService, assetService: assetService}
}

func (dismissalHandler *dismissalHandler) GenerateDismissalDocuments() error {
	csv, err := dismissalHandler.dismissalService.ReadCsvFile("info.csv")
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	var dismissalRecords []*models.DismissalRecord

	for i, row := range csv {
		if i == 0 {
			continue
		}

		record, err := dismissalHandler.dismissalService.CreateDismissalRecord(row)
		if err != nil {
			return fmt.Errorf("failed to create dismissal record: %w", err)
		}

		isc := row[1]

		laptop, err := dismissalHandler.assetService.GetByISC(isc)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		description, err := dismissalHandler.assetService.GetLaptopDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		record.Name = description.Name
		record.InventoryID = description.InventoryID
		record.Serial = description.Serial

		dismissalRecords = append(dismissalRecords, record)

	}

	for _, record := range dismissalRecords {
		templateNames := []string{"commitee", "dismissal", "record"}

		dirPath, err := dismissalHandler.dismissalService.CreateOutputDirectory(record.ISC)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		for _, templateName := range templateNames {
			template, err := dismissalHandler.dismissalService.CreateTemplate(record, templateName)
			if err != nil {
				return fmt.Errorf("failed to create template %v: %w", templateName, err)
			}

			object, err := dismissalHandler.dismissalService.CreateObjectFromTemplate(template)
			if err != nil {
				return fmt.Errorf("failed to create object from template %v: %w", templateName, err)
			}

			file, err := dismissalHandler.dismissalService.CreateFile(dirPath, templateName, "pdf")
			if err != nil {
				return fmt.Errorf("failed to create file")
			}
			defer file.Close()

			err = dismissalHandler.dismissalService.CreatePDF(object, file)
			if err != nil {
				return fmt.Errorf("failed to generate pdf: %w", err)
			}
		}
	}

	return nil
}
