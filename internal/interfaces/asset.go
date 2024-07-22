package interfaces

import (
	"main/internal/models"
)

type AssetService interface {
	GetAll(iql string) (*models.GetObjectRes, error)
	Update(payload *models.Object) (*models.Object, error)
	GetByISC(ISC string) (*models.Object, error)
	GetAttachments(*models.Object) ([]models.Attachment, error)
	DisableUser(user *models.Object) (*models.Object, error)
	SetUserCategory(user *models.Object, category string) (*models.Object, error)
	GetUserEmail(user *models.Object) string
	GetUserLaptops(user *models.Object) (*models.GetObjectRes, error)
	GetUserByEmail(email string) (*models.Object, error)
	GetLaptopDescription(laptop *models.Object) (*models.LaptopDescription, error)
	ExtractInformationResourceIdentifier(value string) (string, error)
}
