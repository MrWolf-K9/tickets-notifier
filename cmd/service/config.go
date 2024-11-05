package service

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey       string   `yaml:"apiKey"`
	TicketsURL   string   `yaml:"ticketsUrl"`
	PhoneNumbers []string `yaml:"phoneNumbers"`
	MessageText  string   `yaml:"messageText"`
	RegexString  string   `yaml:"regexString"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
