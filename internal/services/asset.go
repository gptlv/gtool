package services

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *assetService) GetAll(iql string) (*models.GetObjectRes, error) {
	if iql == "" {
		return nil, fmt.Errorf("empty iql")
	}

	res := &models.GetObjectRes{}

	base := "rest/insight/1.0/iql/objects?iql="
	endpoint := base + iql + "&objectSchemaId=41"

	req, err := s.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return res, nil
}

func (s *assetService) Update(payload *models.Object) (*models.Object, error) {
	if payload == nil {
		return nil, fmt.Errorf("empty payload")
	}

	object := &models.Object{}

	endpoint := fmt.Sprintf("rest/insight/1.0/object/%v/", payload.ObjectKey)

	req, err := s.client.NewRequest("PUT", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = s.client.Do(req, object)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return object, nil
}

func (s *assetService) GetByISC(ISC string) (*models.Object, error) {
	objectEndPoint := fmt.Sprintf(endpoints.GetObjectByISC, ISC)

	req, err := s.client.NewRequest("GET", objectEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	object := new(models.Object)

	_, err = s.client.Do(req, object)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return object, nil

}

func (s *assetService) GetAttachments(object *models.Object) ([]models.Attachment, error) {
	objectAttachmentsEndPoint := fmt.Sprintf(endpoints.GetAttachments, object.ID)

	req, err := s.client.NewRequest("GET", objectAttachmentsEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	attachments := new([]models.Attachment)

	_, err = s.client.Do(req, attachments)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return *attachments, nil

}

func (s *assetService) DisableUser(user *models.Object) (*models.Object, error) {
	payload := new(UserAttributesPayload)
	body := fmt.Sprintf(userAttributePayloadBody, 4220, 100)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	userCategoryEndPoint := fmt.Sprintf(endpoints.GetObjectByISC, user.ID)

	req, err := s.client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = s.client.Do(req, user)
	if err != nil {
		return user, fmt.Errorf("failed to do a request: %w", err)
	}

	return user, nil
}

func (s *assetService) SetUserCategory(user *models.Object, category string) (*models.Object, error) {
	if category != "BYOD" && category != "Corporate laptop" {
		return nil, errors.New("invalid user category")
	}

	payload := new(UserAttributesPayload)
	body := fmt.Sprintf(userAttributePayloadBody, USER_CATEGORY_ATTRIBUTE_ID, "BYOD")

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	userCategoryEndPoint := fmt.Sprintf(endpoints.GetObjectByISC, user.ObjectKey)

	req, err := s.client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = s.client.Do(req, user)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return user, nil
}

func (s *assetService) Search(endpoint string) (*GetObjectRes, error) {
	req, err := s.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	res := new(GetObjectRes)

	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return res, nil
}

func (s *assetService) GetUserByEmail(email string) (*models.Object, error) {
	if email == "" {
		return nil, fmt.Errorf("empty email")
	}
	userEndPoint := fmt.Sprintf(endpoints.GetUserByEmail, email)

	res, err := s.Search(userEndPoint)
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

func (s *assetService) GetUserLaptops(user *models.Object) (*models.GetObjectRes, error) {
	email := s.GetUserEmail(user)
	if email == "" {
		return nil, fmt.Errorf("empty user email")
	}

	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("no such user")
	}

	getUsersLaptopsQuery := fmt.Sprintf("object+having+outboundReferences(Key+=+%v)+and+objectType+=+Laptops", user.ObjectKey)

	laptops, err := s.GetAll(getUsersLaptopsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get %v's laptops", email)
	}

	return laptops, nil
}

func (s *assetService) GetUserEmail(user *models.Object) string {
	for _, attr := range user.Attributes {
		if attr.ObjectTypeAttributeID == USER_EMAIL_ATTRIBUTE_ID {
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
		case NAME_ATTRIBUTE_ID:
			description.Name = attributeValue
		case ISC_ATTRIBUTE_ID:
			description.ISC = attributeValue
		case SERIAL_ATTRIBUTE_ID:
			description.Serial = attributeValue
		case COST_ATTRIBUTE_ID:
			description.Cost = attributeValue
		case INVENTORY_ID_ATTRIBUTE_ID:
			description.InventoryID = attributeValue
		}
	}

	return description, nil
}
