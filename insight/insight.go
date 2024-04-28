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
	"rest/insight/1.0/iql/objects?iql=Email=%v&objectSchemaId=41", // 41 -- IT SD CMDB
	"rest/insight/1.0/object/%v/",
	"rest/insight/1.0/attachments/object/%v",
}

func GetObjectByISC(client *jira.Client, ISC string) (*types.ObjectEntry, error) {
	objectEndPoint := fmt.Sprintf(endpoints.GetObjectByISC, ISC)

	req, err := client.NewRequest("GET", objectEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	object := new(types.ObjectEntry)

	_, err = client.Do(req, object)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return object, nil

}

func GetObjectAttachments(client *jira.Client, object *types.ObjectEntry) ([]types.InsightAttachment, error) {
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

func DisableUser(client *jira.Client, user *types.ObjectEntry) (*types.ObjectEntry, error) {
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

func SetUserCategory(client *jira.Client, user *types.ObjectEntry, category string) (*types.ObjectEntry, error) {
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

func SearchForObjects(client *jira.Client, endpoint string) (*types.ObjectSearchRes, error) {
	req, err := client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	res := new(types.ObjectSearchRes)

	_, err = client.Do(req, res)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	return res, nil
}

func GetUserByEmail(client *jira.Client, email string) (*types.ObjectEntry, error) {
	userEndPoint := fmt.Sprintf(endpoints.GetUserByEmail, email)

	res, err := SearchForObjects(client, userEndPoint)
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

func GetUserLaptops(client *jira.Client, user *types.ObjectEntry) ([]types.ObjectEntry, error) {
	if user == nil {
		return nil, fmt.Errorf("empty user")
	}
	userLaptopsEndPoint := fmt.Sprintf(endpoints.GetUserLaptopsByKey, user.ObjectKey)

	entries, err := SearchForObjects(client, userLaptopsEndPoint)
	if err != nil {
		return nil, fmt.Errorf("failed for search for laptops: %w", err)
	}

	return entries.ObjectEntries, nil
}

func GetUserEmail(user *types.ObjectEntry) string {
	for _, attr := range user.Attributes {
		if attr.ObjectTypeAttributeID == USER_EMAIL_ATTRIBUTE_ID {
			return attr.ObjectAttributeValues[0].Value
		}
	}

	return ""
}
