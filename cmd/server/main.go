package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alukart32/go-fast-key/internal/application"
	"github.com/alukart32/go-fast-key/internal/configuration"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := &configuration.Config{}
	if ConfigFileName != "" {
		data, err := os.ReadFile(ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err = configuration.Load(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	app, err := application.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("App is created")

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
