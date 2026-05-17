package config

import (
	"log"
	"os"

	"gitlab.com/MikeTTh/env"
	"gopkg.in/yaml.v3"
)

type Application struct {
	Name   string   `yaml:"name"`
	Icon   string   `yaml:"icon"`
	Hidden bool     `yaml:"hidden"`
	Exec   []string `yaml:"exec"`
}

type Root struct {
	Applications []Application `yaml:"applications"`
}

var Config Root

func LoadConfig() error {
	configPath := env.String("FRIG_CONFIG", "/etc/frig/config.yaml")

	f, err := os.OpenFile(configPath, os.O_RDONLY, 0)
	if err != nil {
		log.Println("Failed to open config file")
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("WARNING: Failed to close config file:", err)
		}
	}(f)

	return yaml.NewDecoder(f).Decode(&Config)
}
