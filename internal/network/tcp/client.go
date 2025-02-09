package network

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TCPClient struct {
	conn        net.Conn
	idleTimeout time.Duration
	bufferSize  int
}

func NewTcpClient(address string, options ...TCPClientOption) (*TCPClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &TCPClient{
		conn:        conn,
		idleTimeout: 300 * time.Second,
		bufferSize:  4096,
	}

	for _, opt := range options {
		opt(client)
	}

	if client.idleTimeout != 0 {
		if err := client.conn.SetDeadline(time.Now().Add(client.idleTimeout)); err != nil {
			return nil, fmt.Errorf("tcp client: set deadline: %w", err)
		}
	}

	return client, nil
}

func (c *TCPClient) Send(data []byte) ([]byte, error) {
	if _, err := c.conn.Write(data); err != nil {
		return nil, fmt.Errorf("tcp client: write: %w", err)
	}

	response := make([]byte, c.bufferSize)
	count, err := c.conn.Read(response)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("tcp client: read: %w", err)
	} else if count == c.bufferSize {
		return nil, fmt.Errorf("tcp client: response buffer overflow")
	}

	return response[:count], nil
}

func (c *TCPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
