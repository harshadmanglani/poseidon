package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Workflows struct {
		Logs struct {
			Enabled  bool   `yaml:"enabled"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"logs"`
		Metrics struct {
			Enabled  bool   `yaml:"enabled"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"metrics"`
		Incidents struct {
			Enabled  bool   `yaml:"enabled"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"incidents"`
		Operations struct {
			Enabled  bool   `yaml:"enabled"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"operations"`
	} `yaml:"workflows"`

	Clients struct {
		Anthropic struct {
			Key string `yaml:"key"`
		} `yaml:"anthropic"`
	} `yaml:"clients"`
	
	Database struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`
}

var PoseidonConf Config

func Init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(data, &PoseidonConf); err != nil {
		log.Fatal(err)
	}
}
