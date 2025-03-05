package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/alukart32/go-fast-key/internal/pkg/concurrency"
	"go.uber.org/zap"
)

type TCPHandler = func(context.Context, []byte) []byte

type TCPServer struct {
	listener  net.Listener
	semaphore *concurrency.Semaphore

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	logger *zap.Logger
}

func NewTCPServer(address string, logger *zap.Logger, options ...TCPServerOption) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("fail to listen: %w", err)
	}

	server := &TCPServer{
		listener: listener,
		logger:   logger,
	}

	for _, option := range options {
		option(server)
	}

	server.semaphore = concurrency.NewSemaphore(server.maxConnections)
	if server.bufferSize == 0 {
		server.bufferSize = 4 << 10
	}

	return server, nil
}

func (s *TCPServer) HandleQueries(ctx context.Context, handler TCPHandler) {
	if s == nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("fail to accept", zap.Error(err))
				continue
			}

			s.semaphore.Acquire()
			go func(conn net.Conn) {
				defer s.semaphore.Release()
				s.handleConn(ctx, conn, handler)
			}(conn)
		}
	}()

	<-ctx.Done()

	s.listener.Close()
	wg.Wait() // wait goroutine to shut down before all connections are closed.

	return
}

func (s *TCPServer) BufferSize() int {
	return s.bufferSize
}

func (s *TCPServer) MaxConnections() int {
	return s.maxConnections
}

func (s *TCPServer) IdleTimeout() time.Duration {
	return s.idleTimeout
}

func (s *TCPServer) handleConn(ctx context.Context, conn net.Conn, handler TCPHandler) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("captured panic", zap.Any("panic", v))
		}

		if err := conn.Close(); err != nil {
			s.logger.Warn("fail to close connection", zap.Error(err))
		}
	}()

	// reuse buffer for requests
	request := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			if err := conn.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("fail to set read deadline", zap.Error(err))
				break
			}
		}

		count, err := conn.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn(
				"fail to read data",
				zap.String("address", conn.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		} else if count == s.bufferSize {
			s.logger.Warn("small buffer size", zap.Int("buffer_size", s.bufferSize))
			break
		}

		if s.idleTimeout != 0 {
			if err := conn.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("fail to set read deadline", zap.Error(err))
				break
			}
		}

		response := handler(ctx, request[:count])
		if _, err := conn.Write(response); err != nil {
			s.logger.Warn(
				"fail to write data",
				zap.String("address", conn.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		}
	}
}
