package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	appsender "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	internalampq "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/ampq"
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

	ctx := context.Background()

	consumer := internalampq.NewConsumer(config.Ampq)
	if err := consumer.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer consumer.Close(ctx)

	sender := appsender.New(consumer, logger)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go sender.Send(ctx)

	<-ctx.Done()
}
