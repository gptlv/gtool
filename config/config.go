package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

const USER_STATUS_ATTRIBUTE_ID = 4220
const USER_CATEGORY_ATTRIBUTE_ID = 10209
const USER_EMAIL_ATTRIBUTE_ID = 2874
const USER_STATUS_DISABLE_VALUE = 100

const ISC_ATTRIBUTE_ID = 879
const NAME_ATTRIBUTE_ID = 880
const SERIAL_ATTRIBUTE_ID = 889
const COST_ATTRIBUTE_ID = 4184
const INVENTORY_ID_ATTRIBUTE_ID = 932

// type objectService struct {
// 	client *jira.Client
// }

type UserAttributesPayload struct {
	Attributes []struct {
		ObjectTypeAttributeID int `json:"objectTypeAttributeId"`
		ObjectAttributeValues []struct {
			Value string `json:"value"`
		} `json:"objectAttributeValues"`
	} `json:"attributes"`
}

type DismissalRecord struct {
	//comes from csv
	ID       string
	ISC      string
	Flaw     string
	Decision string
	//from insight
	Serial      string
	Name        string
	InventoryID string
	//common
	Date string
	Boss string
	Lead string
}
