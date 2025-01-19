package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/buurzx/in-mem-kvdb/internal/config"
	"github.com/buurzx/in-mem-kvdb/internal/initialization"
	"go.uber.org/zap"
)

var configFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	config := config.MustParseConfiguration(configFileName)

	logger, err := initialization.CreateLogger(
		config.Logger.Level,
		config.Logger.OutputFilePath,
	)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initialization.CreateDatabase(logger)
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("[in-mem-kvdb] > ")

		request, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				logger.Fatal("connection closed", zap.Error(err))
				break
			}

			logger.Error("failed to read from stdin", zap.Error(err))
			continue
		}

		request = strings.TrimSpace(request)
		if request == "exit" {
			logger.Info("exiting in-mem-kvdb")
			break
		}

		db.HandleRequest(request)
	}
}
