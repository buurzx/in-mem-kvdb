package server

import (
	"cmp"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a application configuration structure
type Config struct {
	Engine struct {
		Type string `yaml:"type" env:"ENGINE_TYPE" env-description:"Database engine type" env-default:"in_memory"`
	} `yaml:"engine"`

	Network struct {
		Address        string `yaml:"address" env:"KVDB_ADDRESS" env-description:"Network address"`
		MaxConnections int    `yaml:"max-connections" env:"KVDB_MAX_CONNECTIONS" env-description:"Maximum number of connections" env-default:"100"`
		MaxMessageSize int    `yaml:"max-message-size" env:"KVDB_NETWORK_MAX_MESSAGE_SIZE" env-description:"Maximum message size" env-default:"4096"`
		IdleTimeout    int    `yaml:"idle-timeout" env:"KVDB_NETWORK_IDLE_TIMEOUT" env-description:"Idle timeout" env-default:"300"`
	} `yaml:"network"`

	Logger struct {
		Level          string `yaml:"level" env:"LOG_LEVEL" env-description:"Log level" env-default:"info"`
		OutputFilePath string `yaml:"output-file-path" env:"LOG_OUTPUT_FILE_PATH" env-description:"Log output file path"`
	} `yaml:"logger"`
}

func mustParseConfiguration() *Config {
	const confg = "config.yml"
	var (
		cfg            Config
		configFileName = cmp.Or(os.Getenv("CONFIG_FILE"), confg)
	)

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(configFileName, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
