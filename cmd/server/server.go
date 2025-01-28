package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/buurzx/in-mem-kvdb/internal/initialization"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func BuildCmd() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the in-memory key-value database server",
		Action: func(c *cli.Context) error {
			return runServer()
		},
	}
}

func runServer() error {
	config := mustParseConfiguration()

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

	var ctx context.Context

	for {
		ctx = context.Background()

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

		response := db.HandleRequest(ctx, request)
		fmt.Println(response)
	}

	return nil
}
