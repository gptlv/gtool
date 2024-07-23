package services

import (
	"main/internal/domain"
	"main/internal/interfaces"
	"os"
	"time"

	"github.com/goodsign/monday"
)

type writeOffService struct {
}

func NewWriteOffService() interfaces.WriteOffService {
	return &writeOffService{}
}

func (writeOffService *writeOffService) CreateWriteOffRecord(row []string) (*domain.WriteOffRecord, error) {
	record := new(domain.WriteOffRecord)

	for i, value := range row {
		if i == 0 {
			record.ID = value
		}

		if i == 1 {
			record.ISC = value
		}

		if i == 2 {
			record.Flaw = value
		}

		if i == 3 {
			record.Decision = value
		}
	}

	boss := os.Getenv("BOSS")
	lead := os.Getenv("LEAD")

	t := time.Now()
	layout := "2 January 2006"
	date := monday.Format(t, layout, monday.LocaleRuRU)

	record.Date = date
	record.Boss = boss
	record.Lead = lead

	return record, nil
}

func (writeOffService *writeOffService) TransformDataTo2DSlice(writeOffRecords []*domain.WriteOffRecord) [][]string {
	numRows := len(writeOffRecords)
	result := make([][]string, numRows+1)

	result[0] = []string{"id", "isc", "flaw", "decision", "serial", "name", "inventory_id", "date", "boss", "lead"}

	for i, record := range writeOffRecords {
		result[i+1] = []string{record.ID, record.ISC, record.Flaw, record.Decision, record.Serial, record.Name, record.InventoryID, record.Date, record.Boss, record.Lead}
	}

	return result
}
