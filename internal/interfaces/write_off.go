package interfaces

import "github.com/gptlv/gtools/internal/domain"

type WriteOffService interface {
	CreateWriteOffRecord(row []string) (*domain.WriteOffRecord, error)
	TransformDataTo2DSlice(writeOffRecords []*domain.WriteOffRecord) [][]string
}

type WriteOffHandler interface {
	GenerateWriteOffRecords() error
}
