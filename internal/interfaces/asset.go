package interfaces

import (
	"main/internal/entities"
)

type AssetService interface {
	GetAll(iql string) ([]*entities.Asset, error)
	Update(payload *entities.Asset) (*entities.Asset, error)
	GetByISC(ISC string) (*entities.Asset, error)
	// GetAttachments(*entities.Asset) ([]Attachment, error)
	DisableUser(user *entities.Asset) (*entities.Asset, error)
	SetUserCategory(user *entities.Asset, category string) (*entities.Asset, error)
	GetUserEmail(user *entities.Asset) string
	GetUserLaptops(user *entities.Asset) (*entities.Asset, error)
	GetUserByEmail(email string) (*entities.Asset, error)
	// GetLaptopDescription(laptop *entities.Object) (*AssetDescription, error)
}

type AssetUsecase interface {
	PrintLaptopDescription()
}
