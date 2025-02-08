package kvcli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/buurzx/in-mem-kvdb/internal/initialization"
	network "github.com/buurzx/in-mem-kvdb/internal/network/tcp"
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
			return runCli(c.Context, cfg)
		},
	}
}

func runCli(ctx context.Context, cfg *Config) error {
	ctxWithCancel, cancel := context.WithCancel(ctx)

	logger, err := initialization.CreateLogger(
		cfg.Logger.Level,
		cfg.Logger.OutputFilePath,
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting kvdb-cli")

	reader := bufio.NewReader(os.Stdin)

	options := []network.TCPClientOption{
		network.WithClientIdleTimeout(time.Duration(cfg.Network.IdleTimeout)),
		network.WithBufferSize(cfg.Network.MaxMessageSize),
	}

	client, err := network.NewTcpClient(cfg.Network.Address, options...)
	if err != nil {
		logger.Fatal("failed to create tcp client", zap.Error(err))
	}

	// create tcp connections to the server
	// Channel for handling client requests
	requestChan := make(chan string)
	defer close(requestChan)

	// Start goroutine to handle client requests
	go handleRequests(ctxWithCancel, client, logger, requestChan)
	handleSignals(cancel, client, logger)

	for {
		if err := processRequest(ctxWithCancel, reader, requestChan, cancel, logger); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			logger.Error("process request error", zap.Error(err))
		}
	}
}

func handleRequests(
	ctx context.Context,
	client *network.TCPClient,
	logger *zap.Logger,
	requestChan <-chan string,
) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic in client handler", zap.Any("error", r))
		}
	}()

requestLoop:
	for {
		select {
		case <-ctx.Done():
			break requestLoop
		case request, ok := <-requestChan:
			if !ok {
				logger.Info("request channel closed")
				break requestLoop
			}
			response, err := client.Send([]byte(request))
			if err != nil {
				logger.Error("failed to send request to server", zap.Error(err))
				continue
			}
			fmt.Println(string(response))
		}
	}
}

func handleSignals(cancel context.CancelFunc, client *network.TCPClient, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("received signal", zap.String("signal", sig.String()))

		if err := client.Close(); err != nil {
			logger.Error("error closing client connection", zap.Error(err))
		}

		cancel()
	}()
}

func processRequest(ctx context.Context, reader *bufio.Reader, requestChan chan<- string, cancel context.CancelFunc, logger *zap.Logger) error {
	// Check if context is done before reading
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		fmt.Print("[in-mem-kvdb] > ")
	}

	// Use a channel to handle read timeout
	readChan := make(chan string)
	errChan := make(chan error)

	go func() {
		request, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			return
		}
		readChan <- request
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection closed", zap.Error(err))
			cancel()
			return err
		}
		logger.Error("failed to read from stdin", zap.Error(err))
		return err
	case request := <-readChan:
		request = strings.TrimSpace(request)
		if request == "exit" {
			logger.Info("exiting in-mem-kvdb")
			cancel()
			return nil
		}

		logger.Info("write request to the channel")
		requestChan <- request
		logger.Info("wrote request to the channel - done")
	}

	return nil
}
