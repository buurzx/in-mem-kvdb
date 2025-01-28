package kvcli

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
	cfg := mustParseConfiguration()

	return &cli.Command{
		Name:  "kvdb-cli",
		Usage: "A simple CLI for interacting with the in-memory key-value database",
		Flags: cfg.buildFlags(),
		Action: func(c *cli.Context) error {
			return runCli(c, cfg)
		},
	}
}

func runCli(ctx context.Context, cfg *Config) error {
	logger, err := initialization.CreateLogger(
		cfg.Logger.Level,
		cfg.Logger.OutputFilePath,
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting kvdb-cli")

	reader := bufio.NewReader(os.Stdin)

	// create tcp connections to the server
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

		response := db.HandleRequest(ctx, request)
		fmt.Println(response)
	}

	return nil
}
