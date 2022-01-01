package main

import (
	"context"
	"flag"
	"log"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender.dev.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewSenderConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	_ = context.Background()

	log.Printf("%+v\n", logger)
}
