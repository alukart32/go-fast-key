package network_test

import (
	"testing"
	"time"

	"github.com/alukart32/go-fast-key/internal/network"
	"github.com/stretchr/testify/assert"
)

func TestWithServerIdleTimeout(t *testing.T) {
	t.Parallel()

	idleTimeout := time.Second
	option := network.WithServerIdleTimeout(time.Second)

	var server network.TCPServer
	option(&server)

	assert.Equal(t, idleTimeout, server.IdleTimeout())
}

func TestWithServerBufferSize(t *testing.T) {
	t.Parallel()

	var bufferSize uint = 10 << 10
	option := network.WithServerBufferSize(bufferSize)

	var server network.TCPServer
	option(&server)

	assert.Equal(t, bufferSize, uint(server.BufferSize()))
}

func TestWithServerMaxConnectionsNumber(t *testing.T) {
	t.Parallel()

	var maxConnections uint = 100
	option := network.WithServerMaxConnectionsNumber(maxConnections)

	var server network.TCPServer
	option(&server)

	assert.Equal(t, maxConnections, uint(server.MaxConnections()))
}
