package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Storage       `yaml:"storage"`
	PriceUpdates  time.Duration `yaml:"price_updates"`
	ApiKey        string        `yaml:"api_key"`
	VsCurrency    string        `yaml:"vs_currency"`
	Rest          `yaml:"rest"`
	MaxConcurrent int `yaml:"max_concurrent"`
}
type Storage struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	SslMode  string `yaml:"ssl_mode"`
}

type Rest struct {
	Address string `yaml:"address"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/config.yaml" // Путь по умолчанию
	}

	// Преобразуем относительный путь в абсолютный
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path: %w", err))
	}

	configFile, err := os.Open(absPath)
	if err != nil {
		panic(fmt.Errorf("failed to open config file at %s: %w", absPath, err))
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Errorf("failed to decode config: %w", err))
	}

	return &config
}
