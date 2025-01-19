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
	// Database struct {
	// 	Host        string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
	// 	Port        string `yaml:"port" env:"DB_PORT" env-description:"Database port"`
	// 	Username    string `yaml:"username" env:"DB_USER" env-description:"Database user name"`
	// 	Password    string `env:"DB_PASSWORD" env-description:"Database user password"`
	// 	Name        string `yaml:"db-name" env:"DB_NAME" env-description:"Database name"`
	// 	Connections int    `yaml:"connections" env:"DB_CONNECTIONS" env-description:"Total number of database connections"`
	// } `yaml:"database"`
	// Server struct {
	// 	Host string `yaml:"host" env:"SRV_HOST,HOST" env-description:"Server host" env-default:"localhost"`
	// 	Port string `yaml:"port" env:"SRV_PORT,PORT" env-description:"Server port" env-default:"8080"`
	// } `yaml:"server"`
	// Greeting string `env:"GREETING" env-description:"Greeting phrase" env-default:"Hello!"`
}

func MustParseConfiguration(configFileName string) *Config {
	var cfg Config

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(configFileName, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
