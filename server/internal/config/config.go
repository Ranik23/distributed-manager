package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		HTTP struct {
			Host         string `yaml:"host"`
			Port         string    `yaml:"port"`
			ReadTimeout  string `yaml:"read_timeout"`
			WriteTimeout string `yaml:"write_timeout"`
		} `yaml:"http"`
		GRPC struct {
			Host                  string `yaml:"host"`
			Port                  int    `yaml:"port"`
			MaxConcurrentStreams   int    `yaml:"max_concurrent_streams"`
		} `yaml:"grpc"`
	} `yaml:"server"`
	Database struct {
		Postgres struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			DBName   string `yaml:"dbname"`
			SSLMode  string `yaml:"sslmode"`
		} `yaml:"postgres"`
		BoltDB struct {
			Path     string `yaml:"path"`
			ReadOnly bool   `yaml:"read_only"`
		} `yaml:"boltdb"`
	} `yaml:"database"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}