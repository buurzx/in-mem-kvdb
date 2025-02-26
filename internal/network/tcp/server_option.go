package network

import "time"

type TCPServerOption func(*TCPServer)

func WithServerAddress(address string) TCPServerOption {
	return func(s *TCPServer) {
		s.address = address
	}
}

func WithServerMaxConnections(maxConn int) TCPServerOption {
	return func(s *TCPServer) {
		s.maxConn = maxConn
	}
}

func WithServerIdleTimeout(timeout time.Duration) TCPServerOption {
	return func(s *TCPServer) {
		s.idleTimeout = timeout
	}
}

func WithServerBufferSize(bufferSize int) TCPServerOption {
	return func(s *TCPServer) {
		s.bufferSize = bufferSize
	}
}
