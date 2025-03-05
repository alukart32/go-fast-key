package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/alukart32/go-fast-key/internal/network"
	"github.com/alukart32/go-fast-key/internal/pkg/datasize"
	"go.uber.org/zap"
)

func main() {
	address := flag.String("address", "localhost:8080", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	maxMessageSize, err := datasize.Parse(*maxMessageSizeStr)
	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	var options []network.TCPClientOption
	options = append(options, network.WithClientIdleTimeout(*idleTimeout))
	options = append(options, network.WithClientBufferSize(uint(maxMessageSize)))

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTCPClient(*address, options...)
	if err != nil {
		logger.Fatal("failed to connect with server", zap.Error(err))
	}

	for {
		fmt.Print("[fastkey] > ")
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("fail to read query", zap.Error(err))
		}

		response, err := client.Send([]byte(request))
		logger.Info("response", zap.String("response", string(response)), zap.Error(err))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("fail to send query", zap.Error(err))
		}

		fmt.Println(string(response))

		client.RefreshDeadline()
	}
}
