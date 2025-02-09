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
		network.WithServerIdleTimeout(time.Duration(config.Network.IdleTimeout)*time.Second),
		network.WithServerMaxConnections(config.Network.MaxConnections),
		network.WithServerBufferSize(config.Network.MaxMessageSize),
	)
	if err != nil {
		logger.Fatal("failed to create tcp server", zap.Error(err))
	}

	// Start server in a goroutine
	go func() {
		logger.Info("starting in-mem-kvdb server", zap.String("address", config.Network.Address))
		server.Start(ctxWithCancel, db)
	}()

	handleSignals(cancel, logger)

	// Wait for either context cancellation or server error

	<-ctxWithCancel.Done()
	logger.Info("shutting down server gracefully...")
	if err := server.Close(); err != nil {
		logger.Error("error during shutdown", zap.Error(err))
	}

	os.Exit(0)

	return nil
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
