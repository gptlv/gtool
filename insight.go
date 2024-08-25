package main

import (
	"errors"
	"fmt"

	"github.com/andygrunwald/go-jira"
)

type ObjectDescription struct {
	ISC         string `csv:"isc"`
	Name        string `csv:"name"`
	Cost        string `csv:"cost"`
	Serial      string `csv:"serial"`
	InventoryID string `csv:"inventory_id"`
}

type Record struct {
	ID string `csv:"id,omitempty"`
	*ObjectDescription
	Flaw           string `csv:"flaw"`
	Decision       string `csv:"decision"`
	Date           string `csv:"date,omitempty"`
	TeamLead       string `csv:"team_lead,omitempty"`
	DepartmentLead string `csv:"department_lead,omitempty"`
	Director       string `csv:"director,omitempty"`
}

func (d ObjectDescription) String() string {
	return fmt.Sprintf(`
	Ноутбук: %s 
	Инвентарный номер: %s 
	Серийный номер: %s
	
	Стоимость: %s`, d.Name, d.ISC, d.Serial, d.Cost)
}

func (g *gtool) getUserByEmail(email string) (*jira.ObjectList, *jira.Response, error) {
	findUserPayload := new(jira.FindObjectPayload)

	findUserQuery := fmt.Sprintf("Email == %s", email)
	findUserPayload.Query = findUserQuery
	findUserPayload.ObjectTypeID = "228" // User
	findUserPayload.ResultPerPage = 1
	findUserPayload.ObjectSchemaID = 41 // CMDB
	findUserPayload.IncludeAttributes = true

	return g.client.Object.Find(findUserPayload)
}

func (g *gtool) getUserLaptops(user *jira.Object) (*jira.ObjectList, *jira.Response, error) {
	var email string

	for _, attribute := range user.Attributes {
		if attribute.ObjectTypeAttributeID != 2874 {
			continue
		}

		if len(attribute.ObjectAttributeValues) != 1 {
			return nil, nil, fmt.Errorf("invalid value")
		}
		email = attribute.ObjectAttributeValues[0].Value
	}

	if email == "" {
		return nil, nil, errors.New("empty email")
	}

	findUserLaptopsPayload := new(jira.FindObjectPayload)

	findUserLaptopsQuery := fmt.Sprintf("object having outboundReferences(Email == %s) and objectType == Laptops", email)
	findUserLaptopsPayload.ObjectTypeID = "129" // Laptop
	findUserLaptopsPayload.ResultPerPage = 1
	findUserLaptopsPayload.Query = findUserLaptopsQuery
	findUserLaptopsPayload.ObjectSchemaID = 41
	findUserLaptopsPayload.IncludeAttributes = true

	return g.client.Object.Find(findUserLaptopsPayload)
}

func (g *gtool) disableUser(user *jira.Object) (*jira.Object, *jira.Response, error) {
	updateUserPayload := new(jira.UpdateObjectPayload)
	statusAttribute := new(jira.Attribute)
	statusAttribute.ObjectTypeAttributeID = 4220
	statusAttribute.ObjectAttributeValues = []jira.ObjectAttributeValue{{Value: "100"}}
	updateUserPayload.Attributes = []jira.Attribute{*statusAttribute}

	return g.client.Object.Update(user.ObjectKey, updateUserPayload)
}

func (g *gtool) setUserCategory(user *jira.Object, category string) (*jira.Object, *jira.Response, error) {
	updateUserPayload := new(jira.UpdateObjectPayload)
	categoryAttribute := new(jira.Attribute)
	categoryAttribute.ObjectTypeAttributeID = 10209
	categoryAttribute.ObjectAttributeValues = []jira.ObjectAttributeValue{{Value: category}}
	updateUserPayload.Attributes = []jira.Attribute{*categoryAttribute}

	return g.client.Object.Update(user.ObjectKey, updateUserPayload)
}

func (g *gtool) getObjectDescription(object *jira.Object) (*ObjectDescription, error) {
	attributeMap := map[int]string{
		config.Jira.Attribute.ISC:         "ISC",
		config.Jira.Attribute.Name:        "Name",
		config.Jira.Attribute.Serial:      "Serial",
		config.Jira.Attribute.Cost:        "Cost",
		config.Jira.Attribute.InventoryID: "InventoryID",
	}
	description := new(ObjectDescription)

	for _, attribute := range object.Attributes {
		if fieldName, ok := attributeMap[attribute.ObjectTypeAttributeID]; ok {
			value := attribute.ObjectAttributeValues[0].Value
			switch fieldName {
			case "Name":
				description.Name = value
			case "ISC":
				description.ISC = value
			case "Serial":
				description.Serial = value
			case "Cost":
				description.Cost = value
			case "InventoryID":
				description.InventoryID = value
			}
		}
	}

	return description, nil
}
