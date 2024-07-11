package controllers

import (
	"fmt"
	"main/internal/interfaces"
)

type DismissalHandler struct {
	dismissalService interfaces.DismissalService
	assetService     interfaces.AssetService
}

func NewDismissalHandler(dismissalService interfaces.DismissalService, assetService interfaces.AssetService) *DismissalHandler {
	return &DismissalHandler{dismissalService: dismissalService, assetService: assetService}
}

func (dismissalHandler *DismissalHandler) GenerateDismissalDocuments() error {
	csv, err := h.dismissalService.ReadCsvFile("info.csv")
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	var dismissalRecords []*dismissal.DismissalRecord

	for i, row := range csv {
		if i == 0 {
			continue
		}

		record, err := h.dismissalService.CreateDismissalRecord(row)
		if err != nil {
			return fmt.Errorf("failed to create dismissal record: %w", err)
		}

		isc := row[1]

		laptop, err := h.objectService.GetByISC(isc)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		description, err := h.objectService.GetLaptopDescription(laptop)
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

		dirPath, err := h.dismissalService.CreateOutputDirectory(record.ISC)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		for _, templateName := range templateNames {
			template, err := h.dismissalService.CreateTemplate(record, templateName)
			if err != nil {
				return fmt.Errorf("failed to create template %v: %w", templateName, err)
			}

			object, err := h.dismissalService.CreateObjectFromTemplate(template)
			if err != nil {
				return fmt.Errorf("failed to create object from template %v: %w", templateName, err)
			}

			file, err := h.dismissalService.CreateFile(dirPath, templateName, "pdf")
			if err != nil {
				return fmt.Errorf("failed to create file")
			}
			defer file.Close()

			err = h.dismissalService.CreatePDF(object, file)
			if err != nil {
				return fmt.Errorf("failed to generate pdf: %w", err)
			}
		}
	}

	return nil
}
