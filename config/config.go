package config

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Title      string
	MainGenres []string `yaml:"main_genres,flow"`
	Email      Email
}

type Email struct {
	Enabled    bool
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
	From       EmailRecipient
	To         EmailRecipient
}

type EmailRecipient struct {
	Address string
	Name    string
}

/**
 * Parse the contents of the YAML file into the Config object.
 */
func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, &c)
}

func (c *Config) IsInterestingMainGenre(genres []string) bool {
	for _, g := range genres {
		if isContainedInList(g, c.MainGenres) {
			return true
		}
	}
	return false
}

func isContainedInList(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}
	return false
}
