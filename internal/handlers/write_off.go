package handlers

import (
	"fmt"
	"main/internal/domain"
	"main/internal/interfaces"
	"main/util"
)

type writeOffHandler struct {
	writeOffService interfaces.WriteOffService
	assetService    interfaces.AssetService
}

func NewWriteOffHandler(writeOffService interfaces.WriteOffService, assetService interfaces.AssetService) *writeOffHandler {
	return &writeOffHandler{writeOffService: writeOffService, assetService: assetService}
}

func (writeOffHandler *writeOffHandler) GenerateWriteOffRecords() error {
	outputFileName := "write_off_records.csv"
	csv, err := util.ReadCsvFile("info.csv")
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	var writeOffRecords []*domain.WriteOffRecord

	for i, row := range csv {
		if i == 0 {
			continue
		}

		record, err := writeOffHandler.writeOffService.CreateWriteOffRecord(row)
		if err != nil {
			return fmt.Errorf("failed to create dismissal record: %w", err)
		}

		isc := row[1]

		laptop, err := writeOffHandler.assetService.GetByISC(isc)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		description, err := writeOffHandler.assetService.GetLaptopDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		record.Name = description.Name
		record.InventoryID = description.InventoryID
		record.Serial = description.Serial

		writeOffRecords = append(writeOffRecords, record)

	}

	dataToWrite := writeOffHandler.writeOffService.TransformDataTo2DSlice(writeOffRecords)

	err = util.WriteCsvFile(dataToWrite, outputFileName)
	if err != nil {
		return fmt.Errorf("failed to write csv file: %w", err)
	}

	return nil
}
