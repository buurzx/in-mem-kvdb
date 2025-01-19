package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a application configuration structure
type Config struct {
	Logger struct {
		Level          string `yaml:"level" env:"LOG_LEVEL" env-description:"Log level" env-default:"info"`
		OutputFilePath string `yaml:"output-file-path" env:"LOG_OUTPUT_FILE_PATH" env-description:"Log output file path" env-default:"logs.log"`
	}
}

func MustParseConfiguration(configFileName string) *Config {
	var cfg Config

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(configFileName, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
