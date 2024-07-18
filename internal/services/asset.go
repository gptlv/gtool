package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/config"
	"main/internal/interfaces"
	"main/internal/models"

	"github.com/andygrunwald/go-jira"
)

type assetService struct {
	client *jira.Client
}

func NewAssetService(client *jira.Client) interfaces.AssetService {
	return &assetService{client: client}
}

func (assetService *assetService) GetAll(iql string) (*models.GetObjectRes, error) {
	if iql == "" {
		return nil, fmt.Errorf("empty iql")
	}

	res := &models.GetObjectRes{}

	base := "rest/insight/1.0/iql/objects?iql="
	endpoint := base + iql + "&objectSchemaId=41"

	req, err := assetService.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = assetService.client.Do(req, res)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return res, nil
}

func (assetService *assetService) Update(payload *models.Object) (*models.Object, error) {
	if payload == nil {
		return nil, fmt.Errorf("empty payload")
	}

	object := &models.Object{}

	endpoint := fmt.Sprintf("rest/insight/1.0/object/%v/", payload.ObjectKey)

	req, err := assetService.client.NewRequest("PUT", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = assetService.client.Do(req, object)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return object, nil
}

func (assetService *assetService) GetByISC(ISC string) (*models.Object, error) {
	objectEndPoint := fmt.Sprintf(models.Endpoints.GetObjectByISC, ISC)

	req, err := assetService.client.NewRequest("GET", objectEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	object := new(models.Object)

	_, err = assetService.client.Do(req, object)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return object, nil

}

func (assetService *assetService) GetAttachments(object *models.Object) ([]models.Attachment, error) {
	objectAttachmentsEndPoint := fmt.Sprintf(models.Endpoints.GetAttachments, object.ID)

	req, err := assetService.client.NewRequest("GET", objectAttachmentsEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	attachments := new([]models.Attachment)

	_, err = assetService.client.Do(req, attachments)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return *attachments, nil

}

func (assetService *assetService) DisableUser(user *models.Object) (*models.Object, error) {
	payload := new(models.UserAttributesPayload)
	body := fmt.Sprintf(config.UserAttributePayloadBody, 4220, 100)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	userCategoryEndPoint := fmt.Sprintf(models.Endpoints.GetObjectByISC, user.ID)

	req, err := assetService.client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = assetService.client.Do(req, user)
	if err != nil {
		return user, fmt.Errorf("failed to do a request: %w", err)
	}

	return user, nil
}

func (assetService *assetService) SetUserCategory(user *models.Object, category string) (*models.Object, error) {
	if category != "BYOD" && category != "Corporate laptop" {
		return nil, errors.New("invalid user category")
	}

	payload := new(models.UserAttributesPayload)
	body := fmt.Sprintf(config.UserAttributePayloadBody, config.USER_CATEGORY_ATTRIBUTE_ID, "BYOD")

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	userCategoryEndPoint := fmt.Sprintf(models.Endpoints.GetObjectByISC, user.ObjectKey)

	req, err := assetService.client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = assetService.client.Do(req, user)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return user, nil
}

func (assetService *assetService) Search(endpoint string) (*models.GetObjectRes, error) {
	req, err := assetService.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	res := new(models.GetObjectRes)

	_, err = assetService.client.Do(req, res)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return res, nil
}

func (assetService *assetService) GetUserByEmail(email string) (*models.Object, error) {
	if email == "" {
		return nil, fmt.Errorf("empty email")
	}
	userEndPoint := fmt.Sprintf(models.Endpoints.GetUserByEmail, email)

	res, err := assetService.Search(userEndPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to search for a user: %w", err)
	}

	if len(res.ObjectEntries) == 0 {
		return nil, nil
	}

	if len(res.ObjectEntries) > 1 {
		return &res.ObjectEntries[0], fmt.Errorf("found more than one user")
	}

	return &res.ObjectEntries[0], nil
}

func (assetService *assetService) GetUserLaptops(user *models.Object) (*models.GetObjectRes, error) {
	email := assetService.GetUserEmail(user)
	if email == "" {
		return nil, fmt.Errorf("empty user email")
	}

	user, err := assetService.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("no such user")
	}

	getUsersLaptopsQuery := fmt.Sprintf("object+having+outboundReferences(Key+=+%v)+and+objectType+=+Laptops", user.ObjectKey)

	laptops, err := assetService.GetAll(getUsersLaptopsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get %v's laptops", email)
	}

	return laptops, nil
}

func (assetService *assetService) GetUserEmail(user *models.Object) string {
	for _, attr := range user.Attributes {
		if attr.ObjectTypeAttributeID == config.USER_EMAIL_ATTRIBUTE_ID {
			return attr.ObjectAttributeValues[0].Value
		}
	}

	return ""
}

func (s *assetService) GetLaptopDescription(laptop *models.Object) (*models.LaptopDescription, error) {
	if laptop == nil {
		return nil, fmt.Errorf("empty laptop")
	}

	description := &models.LaptopDescription{}

	for _, attribute := range laptop.Attributes {
		attributeValue := attribute.ObjectAttributeValues[0].Value

		switch attribute.ObjectTypeAttributeID {
		case config.NAME_ATTRIBUTE_ID:
			description.Name = attributeValue
		case config.ISC_ATTRIBUTE_ID:
			description.ISC = attributeValue
		case config.SERIAL_ATTRIBUTE_ID:
			description.Serial = attributeValue
		case config.COST_ATTRIBUTE_ID:
			description.Cost = attributeValue
		case config.INVENTORY_ID_ATTRIBUTE_ID:
			description.InventoryID = attributeValue
		}
	}

	return description, nil
}
