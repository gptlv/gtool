package interfaces

import "main/internal/domain"

type WriteOffService interface {
	CreateWriteOffRecord(row []string) (*domain.WriteOffRecord, error)
	TransformDataTo2DSlice(writeOffRecords []*domain.WriteOffRecord) [][]string
}

type WriteOffHandler interface {
	GenerateWriteOffRecords() error
}
