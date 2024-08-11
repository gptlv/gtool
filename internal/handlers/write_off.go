package handlers

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gptlv/gtools/internal/domain"
	"github.com/gptlv/gtools/internal/interfaces"
	"github.com/gptlv/gtools/util"
)

type writeOffHandler struct {
	writeOffService interfaces.WriteOffService
	assetService    interfaces.AssetService
}

func NewWriteOffHandler(writeOffService interfaces.WriteOffService, assetService interfaces.AssetService) *writeOffHandler {
	return &writeOffHandler{writeOffService: writeOffService, assetService: assetService}
}

func (writeOffHandler *writeOffHandler) GenerateWriteOffRecords() error {
	inputFileName := "info.csv"
	outputFileName := "write_off_records.csv"

	log.Info(fmt.Sprintf("input file name: %v", inputFileName))
	log.Info(fmt.Sprintf("output file name: %v\n", outputFileName))

	log.Info(fmt.Sprintf("reading csv input file %v\n", inputFileName))
	csv, err := util.ReadCsvFile(inputFileName)
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	var writeOffRecords []*domain.WriteOffRecord

	for i, row := range csv {
		log.SetPrefix("")
		if i == 0 {
			continue
		}

		log.Info(fmt.Sprintf("generating write-off record for row %v", i))
		record, err := writeOffHandler.writeOffService.CreateWriteOffRecord(row)
		if err != nil {
			return fmt.Errorf("failed to create write-off record: %w", err)
		}

		isc := row[1]
		log.SetPrefix(isc)

		log.Info("getting laptop by isc")
		laptop, err := writeOffHandler.assetService.GetByISC(isc)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		log.Info("getting laptop description")
		description, err := writeOffHandler.assetService.GetLaptopDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		record.Name = description.Name
		record.InventoryID = description.InventoryID
		record.Serial = description.Serial

		log.Info("appending write-off record\n")
		writeOffRecords = append(writeOffRecords, record)

	}

	log.SetPrefix("")
	log.Info("transforming data to 2D slice")
	dataToWrite := writeOffHandler.writeOffService.TransformDataTo2DSlice(writeOffRecords)

	log.Info(fmt.Sprintf("writing data to output file %v", outputFileName))
	err = util.WriteCsvFile(dataToWrite, outputFileName)
	if err != nil {
		return fmt.Errorf("failed to write csv file: %w", err)
	}

	log.Info("finished generating write-off records\n")

	return nil
}
