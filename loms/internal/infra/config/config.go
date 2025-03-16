package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// nolint:revive
// HttpPort, GrpcPort - норма.
type Config struct {
	Service struct {
		Host     string `yaml:"host"`
		HttpPort int    `yaml:"http_port"`
		GrpcPort int    `yaml:"grpc_port"`
		Workers  int    `yaml:"workers"`
	} `yaml:"service"`
	Jaeger struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"jaeger"`
	DbMaster struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
	} `yaml:"db_master"`
	DbReplica struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
	} `yaml:"db_replica"`
	Kafka struct {
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		OrderTopic string `yaml:"order_topic"`
		Brokers    string `yaml:"brokers"`
	} `yaml:"kafka"`
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
