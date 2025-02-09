package kvcli

import (
	"cmp"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/urfave/cli/v2"
)

// Config is a application configuration structure
type Config struct {
	Network struct {
		Address string `yaml:"address" env:"KVDB_ADDRESS" env-description:"Network address"`
		MaxMessageSize int    `yaml:"max-message-size" env:"NETWORK_MAX_MESSAGE_SIZE" env-description:"Maximum message size"`
		IdleTimeout    int    `yaml:"idle-timeout" env:"NETWORK_IDLE_TIMEOUT" env-description:"Idle timeout"`
	} `yaml:"network"`

	Logger struct {
		Level          string `yaml:"level" env:"LOG_LEVEL" env-description:"Log level"`
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

func (c *Config) buildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "kvdb-address",
			EnvVars:     []string{"KVDB_ADDRESS"},
			Value:       c.Network.Address,
			Usage:       "Address of the in-memory key-value database server",
			Destination: &c.Network.Address,
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "logger-level",
			EnvVars:     []string{"LOG_LEVEL"},
			Value:       c.Logger.Level,
			Destination: &c.Logger.Level,
			Usage:       "Log level",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "logger-output-file-path",
			EnvVars:     []string{"LOG_OUTPUT_FILE_PATH"},
			Value:       c.Logger.OutputFilePath,
			Destination: &c.Logger.OutputFilePath,
			Usage:       "Log output file path",
			Required:    false,
		},
	}
}
