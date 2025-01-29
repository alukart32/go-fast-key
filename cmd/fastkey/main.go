package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alukart32/go-fast-key/internal/fastkey"
	"github.com/alukart32/go-fast-key/internal/fastkey/compute"
	"github.com/alukart32/go-fast-key/internal/fastkey/storage"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	storage := storage.NewEngine(256)
	requestParser, err := compute.NewParser(logger)
	if err != nil {
		logger.Fatal("fail to init parser", zap.Error(err))
	}

	db, err := fastkey.NewDatabase(requestParser, storage, logger)
	if err != nil {
		logger.Fatal("fail to init database", zap.Error(err))
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("[fastkey] > ")
		request, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("fail to read request", zap.Error(err))
		}
		if request == "exit" {
			break
		}

		result, err := db.HandleRequest(request)
		if err != nil {
			logger.Error("fail to handle the request", zap.String("request", request), zap.Error(err))
		}
		fmt.Println(result)
	}
}
