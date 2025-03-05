package network_test

import (
	"errors"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/alukart32/go-fast-key/internal/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPClient(t *testing.T) {
	t.Parallel()

	const serverResponse = "hello client"
	const serverAddress = "localhost:11111"

	listener, err := net.Listen("tcp", serverAddress)
	require.NoError(t, err)

	// tcp server for tests.
	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}

			_, err = connection.Read(make([]byte, 2048))
			require.NoError(t, err)

			_, err = connection.Write([]byte(serverResponse))
			require.NoError(t, err)
		}
	}()

	tests := map[string]struct {
		request string
		client  func() *network.TCPClient

		wantResponse string
		wantErr      error
	}{
		"client with incorrect server address": {
			request: "hello server",
			client: func() *network.TCPClient {
				client, err := network.NewTCPClient("localhost:1010")

				var errNo syscall.Errno
				require.ErrorAs(t, err, &errNo)
				require.True(t, errNo == 10061 || errNo == syscall.ECONNREFUSED)
				return client
			},
			wantResponse: serverResponse,
		},
		"client with small max message size": {
			request: "hello server",
			client: func() *network.TCPClient {
				client, err := network.NewTCPClient(serverAddress, network.WithClientBufferSize(5))
				require.NoError(t, err)
				return client
			},
			wantErr: errors.New("small buffer size"),
		},
		"client with idle timeout": {
			request: "hello server",
			client: func() *network.TCPClient {
				client, err := network.NewTCPClient(serverAddress, network.WithClientIdleTimeout(100*time.Millisecond))
				require.NoError(t, err)
				return client
			},
			wantResponse: serverResponse,
		},
		"client without options": {
			request: "hello server",
			client: func() *network.TCPClient {
				client, err := network.NewTCPClient(serverAddress)
				require.NoError(t, err)
				return client
			},
			wantResponse: serverResponse,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := test.client()
			if client == nil {
				return
			}

			response, err := client.Send([]byte(test.request))
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResponse, string(response))
			client.Close()
		})
	}
}
