package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Jira     `yaml:"jira"`
		LDAP     `yaml:"ldap"`
		WriteOff `yaml:"write_off"`
	}

	Jira struct {
		Token     string `env-required:"true" env:"JIRA_TOKEN"`
		URL       string `env-required:"true" env:"JIRA_URL"`
		Attribute `env-required:"true" yaml:"attribute"`
	}

	Attribute struct {
		ISC         int `env-required:"true" yaml:"isc"`
		Name        int `env-required:"true" yaml:"name"`
		Serial      int `env-required:"true" yaml:"serial"`
		Cost        int `env-required:"true" yaml:"cost"`
		InventoryID int `env-required:"true" yaml:"inventory_id"`
	}

	LDAP struct {
		URL           string `env-required:"true" env:"LDAP_URL"`
		BaseDN        string `env-required:"true" env:"LDAP_BASE_DN"`
		AdminDN       string `env-required:"true" env:"ADMIN_DN"`
		AdminPassword string `env-required:"true" env:"ADMIN_PASS"`
	}

	WriteOff struct {
		InputFile      string `env-required:"true" yaml:"input_file" env:"INPUT_FILE"`
		OutputFile     string `env-required:"true" yaml:"output_file" env:"OUTPUT_FILE"`
		DepartmentLead string `env-required:"true" yaml:"department_lead" env:"DEPARTMENT_LEAD"`
		TeamLead       string `env-required:"true" yaml:"team_lead" env:"TEAM_LEAD"`
		Director       string `env-required:"true" yaml:"director" env:"DIRECTOR"`
	}
)

func NewConfig() (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig("./config.yml", config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return config, nil
}
