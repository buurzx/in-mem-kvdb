package network

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/buurzx/in-mem-kvdb/internal/database"
	"go.uber.org/zap"
)

type TCPServer struct {
	listener net.Listener

	address     string
	maxConn     int
	idleTimeout time.Duration
	bufferSize  int

	logger *zap.Logger
}

func (t *TCPServer) Start(ctx context.Context, db *database.Database) error {
	var activeConnections atomic.Int32

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			t.logger.Error("failed to accept connection", zap.Error(err))
			return err
		}

		if int(activeConnections.Load()) >= t.maxConn {
			t.logger.Warn("max connections reached, rejecting connection")
			conn.Close()
			continue
		}

		activeConnections.Add(1)
		t.logger.Info("new connection accepted", zap.Int("active_connections", int(activeConnections.Load())))

		go func(c net.Conn) {
			defer func() {
				c.Close()
				activeConnections.Add(-1)
				t.logger.Info("connection closed", zap.Int("active_connections", int(activeConnections.Load())))
			}()

			buffer := make([]byte, t.bufferSize)
			n, err := c.Read(buffer)
			if err != nil {
				t.logger.Error("failed to read from connection", zap.Error(err))
				return
			}
			request := buffer[:n]

			handler := NewConnectionHandler(c, db, t.logger, t.idleTimeout, t.bufferSize)

			response := handler.Handle(ctx, request)

			_, err = c.Write([]byte(response))
			if err != nil {
				t.logger.Sugar().Errorf("server write to connection %w", err)
			}
		}(conn)
	}
}

func (t *TCPServer) Close() error {
	return t.listener.Close()
}

func NewTCPServer(logger *zap.Logger, options ...TCPServerOption) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("tcp server: logger is required")
	}

	server := &TCPServer{
		address:     "localhost:8080",
		maxConn:     100,
		idleTimeout: 300 * time.Second,
		logger:      logger,
	}

	for _, opt := range options {
		opt(server)
	}

	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		return nil, err
	}

	server.listener = listener

	if server.bufferSize == 0 {
		server.bufferSize = 4 << 10
	}

	return server, nil
}
