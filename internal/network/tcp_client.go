package network

import (
	"errors"
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

func NewTCPClient(address string, options ...TCPClientOption) (*TCPClient, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("fail to dial: %w", err)
	}

	client := &TCPClient{
		conn:       connection,
		bufferSize: defaultBufferSize,
	}

	for _, option := range options {
		option(client)
	}

	if client.idleTimeout != 0 {
		if err := connection.SetDeadline(time.Now().Add(client.idleTimeout)); err != nil {
			return nil, fmt.Errorf("fail to set deadline for connection: %w", err)
		}
	}

	return client, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if _, err := c.conn.Write(request); err != nil {
		return nil, err
	}

	response := make([]byte, c.bufferSize)
	count, err := c.conn.Read(response)
	if err != nil && err != io.EOF {
		return nil, err
	} else if count == c.bufferSize {
		return nil, errors.New("small buffer size")
	}

	return response[:count], nil
}

func (c *TCPClient) BufferSize() int {
	return c.bufferSize
}

func (c *TCPClient) IdleTimeout() time.Duration {
	return c.idleTimeout
}

func (c *TCPClient) RefreshDeadline() error {
	if err := c.conn.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
		return fmt.Errorf("fail to set deadline for connection: %w", err)
	}
	return nil
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
