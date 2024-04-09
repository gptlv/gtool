package insight

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"

	"github.com/andygrunwald/go-jira"
)

// type laptop = struct {
// 	ID        int    `json:"id"`
// 	Label     string `json:"label"`
// 	ObjectKey string `json:"objectKey"`
// }

func GetUserLaptops(client *jira.Client, email string) (*types.InsightObject, error) {
	// type AutoGenerated struct {
	// 	ObjectEntries []laptop `json:"objectEntries"`
	// }

	userLaptopsEndPoint := fmt.Sprintf("rest/insight/1.0/iql/objects?iql=object+having+outboundReferences(Email+=+%v)+and+objectType+=+Laptops", email)

	req, err := client.NewRequest("GET", userLaptopsEndPoint, nil)
	if err != nil {
		return nil, err
	}

	entries := new(types.InsightObject)

	_, err = client.Do(req, entries)
	if err != nil {
		return nil, err
	}

	return entries, nil

}

func DisableUser(client *jira.Client, userISC string) (*types.InsightUser, error) {
	user := new(types.InsightUser)
	if userISC == "" {
		return user, errors.New("empty ISC")
	}

	body := `{
		"attributes": [
		{
			"objectTypeAttributeId": 4220,
			"objectAttributeValues": [
				{
					"value": "100"
				}
			]
		}
		]
	}`

	payload := new(types.InsightUserAttributesPayload)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return user, err
	}

	userCategoryEndPoint := fmt.Sprintf("rest/insight/1.0/object/%v/", userISC)

	req, err := client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return user, err
	}

	_, err = client.Do(req, user)
	if err != nil {
		return user, err
	}

	return user, nil

}

func SetUserCategory(client *jira.Client, userISC string, category string) (*types.InsightUser, error) {
	user := new(types.InsightUser)

	if category != "BYOD" && category != "Corporate laptop" {
		return user, errors.New("invalid user category")
	}

	body := `{
		"attributes": [
		{
			"objectTypeAttributeId": 10209,
			"objectAttributeValues": [
				{
					"value": "BYOD"
				}
			]
		}
		]
	}`

	payload := new(types.InsightUserAttributesPayload)

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return user, err
	}

	userCategoryEndPoint := fmt.Sprintf("rest/insight/1.0/object/%v/", userISC)

	req, err := client.NewRequest("PUT", userCategoryEndPoint, payload)
	if err != nil {
		return user, err
	}

	_, err = client.Do(req, user)
	if err != nil {
		return user, err
	}

	return user, nil

}

func GetUserISC(client *jira.Client, email string) (string, error) {
	userEndPoint := fmt.Sprintf("rest/insight/1.0/iql/objects?iql=Email=%v", email)

	req, err := client.NewRequest("GET", userEndPoint, nil)
	if err != nil {
		return "", err
	}

	user := new(types.InsightUser)

	_, err = client.Do(req, user)
	if err != nil {
		return "", err
	}

	ISC := user.ObjectEntries[0].ObjectKey

	return ISC, nil

}
