package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Jira            `yaml:"jira"`
		ActiveDirectory `yaml:"active_directory"`
		LDAP            `yaml:"ldap"`
		WriteOff        `yaml:"write_off"`
	}

	Jira struct {
		Token string `env-required:"true" env:"JIRA_TOKEN"`
		URL   string `env-required:"true" env:"JIRA_URL"`
	}

	ActiveDirectory struct {
		AdminDN       string `env-required:"true" env:"ADMIN_DN"`
		AdminPassword string `env-required:"true" env:"ADMIN_PASS"`
	}

	LDAP struct {
		URL    string `env-required:"true" env:"LDAP_URL"`
		BaseDN string `env-required:"true" env:"LDAP_BASE_DN"`
	}

	WriteOff struct {
		Boss string `env-required:"true" yaml:"boss" env:"BOSS"`
		Lead string `env-required:"true" yaml:"lead" env:"LEAD"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if envErr := godotenv.Load(".env"); envErr != nil {
		fmt.Println(".env file missing")
	}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
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

const EMAIL_FIELD_KEY = "customfield_10145"

var UserAttributePayloadBody = `{
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

var InternalCommentPayloadBody = `{
	"body": "%s",
	"properties": [
	  {
		"key": "sd.public.comment",
		"value": {
		   "internal": true
		}
	  }
	]
 }`

var BlockByIssuePayloadBody = `
{
	"transition": {
		"id": "%v"
	},
	"update": {
		"issuelinks": [
			{
				"add": {
					"type": {
						"name": "Blocks"
					},
					"inwardIssue": {
						"key": "%v"
					}
				}
			}
		]
	}
}`

var BlockUntilTomorrowPayloadBody = `
{
    "transition": {
        "id": "%v"
    },
    "fields": {
        "customfield_10253": "%v"
    }
}
`
