package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/scheduler/notifier"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	internalampq "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/ampq"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/logger"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler.dev.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewSchedulerConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	storage, err := storage.New(config.Storage)
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Connect(ctx); err != nil {
		logger.Error(err.Error())
		return
	}
	defer storage.Close(ctx)

	producer := internalampq.NewProducer(config.Ampq)
	if err := producer.Connect(ctx); err != nil {
		logger.Error(err.Error())
		return
	}
	defer producer.Close(ctx)

	scheduler := notifier.New(storage, producer, logger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if err := scheduler.NotifyEvents(ctx, 5*time.Minute); err != nil {
		logger.Error(err.Error())
		return
	}
}
