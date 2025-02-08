package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buurzx/in-mem-kvdb/internal/initialization"
	network "github.com/buurzx/in-mem-kvdb/internal/network/tcp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func BuildCmd() *cli.Command {
	return &cli.Command{
		Name:  "kvdb-server",
		Usage: "Start the in-memory key-value database server",
		Action: func(c *cli.Context) error {
			return runServer(c.Context)
		},
	}
}

func runServer(ctx context.Context) error {
	ctxWithCancel, cancel := context.WithCancel(ctx)
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

	server, err := network.NewTCPServer(
		logger,
		network.WithServerAddress(config.Network.Address),
		network.WithServerIdleTimeout(time.Duration(config.Network.IdleTimeout)),
		network.WithServerMaxConnections(config.Network.MaxConnections),
		network.WithServerBufferSize(config.Network.MaxMessageSize),
	)
	if err != nil {
		logger.Fatal("failed to create tcp server", zap.Error(err))
	}

	// Create error channel to handle server errors
	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		logger.Info("starting in-mem-kvdb server", zap.String("address", config.Network.Address))
		if err := server.Start(ctxWithCancel, db); err != nil {
			errChan <- err
		}
	}()

	handleSignals(cancel, logger)

	// Wait for either context cancellation or server error
	select {
	case <-ctxWithCancel.Done():
		logger.Info("shutting down server gracefully...")
		if err := server.Close(); err != nil {
			logger.Error("error during shutdown", zap.Error(err))
		}
		os.Exit(0)
		return nil
	case err := <-errChan:
		return err
	}
}

func handleSignals(cancel context.CancelFunc, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("received signal", zap.String("signal", sig.String()))

		cancel()
	}()
}
