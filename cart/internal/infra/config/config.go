package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Workers int    `yaml:"workers"`
	} `yaml:"service"`
	Jaeger struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"jaeger"`
	ProductService struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Token string `yaml:"token"`
		Limit int    `yaml:"limit"`
		Burst int    `yaml:"burst"`
	} `yaml:"product_service"`
	LomsService struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"loms_service"`
}

func ReadConfig() (Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	configFile = filepath.Clean(configFile)
	if configFile == "" {
		return Config{}, fmt.Errorf("cannot find configFile path")
	}

	log.Println("loading config file", configFile)
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
