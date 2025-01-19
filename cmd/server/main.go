package main

import (
	"log"
	"os"

	"github.com/buurzx/in-mem-kvdb/internal/config"
	"github.com/buurzx/in-mem-kvdb/internal/initialization"
)

var configFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	config := config.MustParseConfiguration(configFileName)

	initializer, err := initialization.NewInitializer(config)
	if err != nil {
		log.Fatal(err)
	}

	initializer.Start()
}
