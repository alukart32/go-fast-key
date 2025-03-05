package network_test

import (
	"testing"
	"time"

	"github.com/alukart32/go-fast-key/internal/network"
	"github.com/stretchr/testify/assert"
)

func TestWithClientIdleTimeout(t *testing.T) {
	t.Parallel()

	idleTimeout := time.Second
	option := network.WithClientIdleTimeout(time.Second)

	var client network.TCPClient
	option(&client)

	assert.Equal(t, idleTimeout, client.IdleTimeout())
}

func TestWithClientBufferSize(t *testing.T) {
	t.Parallel()

	var bufferSize uint = 10 << 10
	option := network.WithClientBufferSize(bufferSize)

	var client network.TCPClient
	option(&client)

	assert.Equal(t, bufferSize, uint(client.BufferSize()))
}
