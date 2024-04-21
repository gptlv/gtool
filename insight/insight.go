package insight

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"

	"github.com/andygrunwald/go-jira"
)

const USER_STATUS_ATTRIBUTE_ID = 4220
const USER_CATEGORY_ATTRIBUTE_ID = 10209
const USER_EMAIL_ATTRIBUTE_ID = 2874
const USER_STATUS_DISABLE_VALUE = 100

var userAttributePayloadBody = `{
	"attributes": [
	{
		"objectTypeAttributeId": %v,
		"objectAttributeValues": [
			{
				"value": "%v"
			}
		]
	}
	]
}`

var endpoints = struct {
	GetUserLaptopsByKey string
	GetUserByEmail      string
	GetObjectByISC      string
	GetAttachments      string
}{
	"rest/insight/1.0/iql/objects?iql=object+having+outboundReferences(Key+=+%v)+and+objectType+=+Laptops",
	"rest/insight/1.0/iql/objects?iql=Email=%v",
	"rest/insight/1.0/object/%v/",
	"rest/insight/1.0/attachments/object/%v",
}

func GetObjectAttachments(client *jira.Client, object *types.InsightObject) ([]types.InsightAttachment, error) {
	objectAttachmentsEndPoint := fmt.Sprintf(endpoints.GetAttachments, object.ID)

	req, err := client.NewRequest("GET", objectAttachmentsEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	attachments := new([]types.InsightAttachment)

	_, err = client.Do(req, attachments)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return *attachments, nil

}

func GetUserLaptops(client *jira.Client, user *types.InsightObject) (*types.InsightObjectEntries, error) {
	userLaptopsEndPoint := fmt.Sprintf(endpoints.GetUserLaptopsByKey, user.ObjectKey)

	req, err := client.NewRequest("GET", userLaptopsEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	entries := new(types.InsightObjectEntries)

	_, err = client.Do(req, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return entries, nil
}

func DisableUser(client *jira.Client, user *types.InsightObject) (*types.InsightObject, error) {
	payload := new(types.InsightUserAttributesPayload)
	body := fmt.Sprintf(userAttributePayloadBody, 4220, 100)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	userCategoryEndPoint := fmt.Sprintf(endpoints.GetObjectByISC, user.ID)

	req, err := client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = client.Do(req, user)
	if err != nil {
		return user, fmt.Errorf("failed to do a request: %w", err)
	}

	return user, nil
}

func SetUserCategory(client *jira.Client, user *types.InsightObject, category string) (*types.InsightObject, error) {
	if category != "BYOD" && category != "Corporate laptop" {
		return nil, errors.New("invalid user category")
	}

	payload := new(types.InsightUserAttributesPayload)
	body := fmt.Sprintf(userAttributePayloadBody, USER_CATEGORY_ATTRIBUTE_ID, "BYOD")

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	userCategoryEndPoint := fmt.Sprintf(endpoints.GetObjectByISC, user.ObjectKey)

	req, err := client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	_, err = client.Do(req, user)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return user, nil
}

func GetUserByEmail(client *jira.Client, email string) (*types.InsightObject, error) {
	userEndPoint := fmt.Sprintf(endpoints.GetUserByEmail, email)

	req, err := client.NewRequest("GET", userEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	entries := new(types.InsightObjectEntries)

	_, err = client.Do(req, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return &entries.ObjectEntries[0], nil
}

func GetUserEmail(user *types.InsightObject) string {
	for _, attr := range user.Attributes {
		if attr.ObjectTypeAttributeID == USER_EMAIL_ATTRIBUTE_ID {
			return attr.ObjectAttributeValues[0].Value
		}
	}

	return ""
}
