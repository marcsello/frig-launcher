package config

import (
	"log"
	"os"

	"github.com/creasty/defaults"
	"gitlab.com/MikeTTh/env"
	"gopkg.in/yaml.v3"
)

const ConfigDefaultPath = "/etc/frig/config.yaml"

type Application struct {
	Name   string   `yaml:"name"`
	Icon   string   `yaml:"icon"`
	Hidden bool     `yaml:"hidden"`
	Exec   []string `yaml:"exec"`
}

type LauncherConfig struct {
	Mode string `yaml:"mode" default:"detached"`
}

type Root struct {
	Launcher     LauncherConfig `yaml:"launcher"`
	Applications []Application  `yaml:"applications"`
}

var Config Root

func LoadConfig() error {

	// set defaults
	err := defaults.Set(&Config)
	if err != nil {
		log.Println("Failed to set defaults:", err)
		return err
	}

	configPath := env.String("FRIG_CONFIG", ConfigDefaultPath)

	f, err := os.OpenFile(configPath, os.O_RDONLY, 0)
	if err != nil {
		log.Println("Failed to open config file", err)
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
