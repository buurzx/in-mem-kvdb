package network

import "time"

type TCPClientOption func(*TCPClient)

func WithClientIdleTimeout(timeout time.Duration) TCPClientOption {
	return func(c *TCPClient) {
		c.idleTimeout = timeout
	}
}

func WithBufferSize(size int) TCPClientOption {
	return func(c *TCPClient) {
		c.bufferSize = size
	}
}
