package network

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buurzx/in-mem-kvdb/internal/database"
	"go.uber.org/zap"
)

type TCPServer struct {
	listener net.Listener

	address           string
	maxConn           int
	idleTimeout       time.Duration
	bufferSize        int
	activeConnections atomic.Int32

	logger *zap.Logger
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

func (t *TCPServer) Start(ctx context.Context, db *database.Database) {
	t.handleQueries(ctx, func(ctx context.Context, request []byte) string {
		response := db.HandleRequest(ctx, string(request))
		return response
	})
}

func (t *TCPServer) Close() error {
	if t.listener == nil {
		return nil
	}

	err := t.listener.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}
	return nil
}

func (t *TCPServer) handleQueries(ctx context.Context, handler func(context.Context, []byte) string) {
	var (
		wg sync.WaitGroup
	)

	wg.Add(1)

	go func() {
		defer wg.Done()

		// reuse buffer
		buffer := make([]byte, t.bufferSize)

		for {
			conn, err := t.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				t.logger.Error("tcp server accept", zap.Error(err))
				continue
			}

			if int(t.activeConnections.Load()) >= t.maxConn {
				t.logger.Error("max connection reached")
				return
			}

			connAmount := t.activeConnections.Add(1)
			t.logger.Info("new connection accepted", zap.Int("active_connections", int(connAmount)))

			t.handleConn(ctx, conn, buffer, handler)
		}

	}()

	<-ctx.Done()
	// Close listener first to stop accepting new connections
	t.Close()
	// Wait for existing connections to finish
	wg.Wait()
}

func (t *TCPServer) handleConn(
	ctx context.Context,
	c net.Conn,
	buffer []byte,
	handler func(context.Context, []byte) string,
) {
	defer func() {
		c.Close()
		val := t.activeConnections.Add(-1)
		t.logger.Info("connection closed", zap.Int("active_connections", int(val)))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if t.idleTimeout > 0 {
				deadline := time.Now().Add(t.idleTimeout)
				t.logger.Info("setting read deadline",
					zap.Duration("timeout", t.idleTimeout),
					zap.Time("deadline", deadline))

				if err := c.SetReadDeadline(deadline); err != nil {
					t.logger.Error("failed to set read deadline", zap.Error(err))
					return
				}
			}

			n, err := c.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					t.logger.Info("connection was closed")
					return
				}

				if errors.Is(err, os.ErrDeadlineExceeded) {
					t.logger.Info("read timed out due to idle timeout", zap.Error(err))
					return
				}

				t.logger.Error("failed to read from connection", zap.Error(err))
				return
			}
			request := buffer[:n]

			if t.idleTimeout > 0 {
				if err := c.SetWriteDeadline(time.Now().Add(t.idleTimeout)); err != nil {
					t.logger.Error("failed to set write deadline", zap.Error(err))
					return
				}
			}

			response := handler(ctx, request)

			_, err = c.Write([]byte(response))
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					t.logger.Info("write timed out due to idle timeout", zap.Error(err))
				} else {
					t.logger.Error("failed to write to connection", zap.Error(err))
				}
				return
			}
		}
	}
}
