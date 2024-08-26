package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
	"github.com/gocarina/gocsv"
	"github.com/goodsign/monday"
)

type ObjectDescription struct {
	ISC         string `csv:"isc"`
	Name        string `csv:"name"`
	Cost        string `csv:"cost"`
	Serial      string `csv:"serial"`
	InventoryID string `csv:"inventory_id"`
}

type Record struct {
	ID int `csv:"id,omitempty"`
	*ObjectDescription
	Flaw           string `csv:"flaw"`
	Decision       string `csv:"decision"`
	Date           string `csv:"date,omitempty"`
	TeamLead       string `csv:"team_lead,omitempty"`
	DepartmentLead string `csv:"department_lead,omitempty"`
	Director       string `csv:"director,omitempty"`
}

func (g *gtool) PrintDescription(email string) error {
	if email == "" {
		return errors.New("empty email")
	}

	userList, _, err := g.getUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(userList.ObjectEntries) == 0 {
		return fmt.Errorf("no user found")
	}

	user := &userList.ObjectEntries[0]

	laptopsRes, _, err := g.getUserLaptops(user)
	if err != nil {
		return fmt.Errorf("failed to get user laptops: %w", err)
	}

	for _, laptop := range laptopsRes.ObjectEntries {
		description, err := g.getObjectDescription(&laptop)
		if err != nil {
			return fmt.Errorf("failed to get description for %s: %w", laptop.ObjectKey, err)
		}

		fmt.Print(description)
	}

	return nil
}

func (g *gtool) GenerateRecords(startID int) error {
	records := []*Record{}

	inputFile := config.WriteOff.InputFile
	outputFile := config.WriteOff.OutputFile

	log.Info(fmt.Sprintf("input file name: %v", inputFile))
	log.Info(fmt.Sprintf("output file name: %v\n", outputFile))

	log.Info(fmt.Sprintf("reading csv input file %v\n", inputFile))
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	log.Info("unmarshaling input file\n")
	err = gocsv.UnmarshalFile(input, &records)
	if err != nil {
		return fmt.Errorf("failed to unmarshal input file %s: %w", inputFile, err)
	}

	t := time.Now()
	layout := "2 January 2006"
	date := monday.Format(t, layout, monday.LocaleRuRU)

	for _, record := range records {
		log.SetPrefix(record.ISC)
		log.Info("getting laptop by isc")
		laptop, _, err := g.client.Object.Get(record.ISC, nil)
		if err != nil {
			return fmt.Errorf("failed to get laptop by ISC: %w", err)
		}

		log.Info("getting laptop description")
		description, err := g.getObjectDescription(laptop)
		if err != nil {
			return fmt.Errorf("failed to get laptop description: %w", err)
		}

		record.ID = startID
		record.ObjectDescription = description
		record.TeamLead = config.WriteOff.TeamLead
		record.DepartmentLead = config.WriteOff.DepartmentLead
		record.Director = config.WriteOff.Director
		record.Date = date

		startID += 1

		log.Info("finished processing the asset\n")
		log.SetPrefix("")
	}

	log.Info("opening output file")
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	log.Info("marshaling output file")
	err = gocsv.MarshalFile(&records, output)
	if err != nil {
		return fmt.Errorf("failed to unmarshal csv file: %w", err)
	}

	log.Info("finished generating write-off records\n")

	return nil
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
