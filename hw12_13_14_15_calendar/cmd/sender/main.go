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

	cfg, err := config.NewSenderConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	ampqTestConf := config.AmpqConf{
		URL:          cfg.Ampq.URL,
		ExchangeType: "fanout",
		ExchangeName: "test",
	}
	producer := internalampq.NewProducer(ampqTestConf)
	if err := producer.Connect(ctx); err != nil {
		logger.Error(err.Error())
		return
	}
	defer producer.Close(ctx)

	consumer := internalampq.NewConsumer(cfg.Ampq)
	if err := consumer.Connect(ctx); err != nil {
		logger.Error(err.Error())
		return
	}
	defer consumer.Close(ctx)

	sender := appsender.New(consumer, producer, logger)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go sender.Send(ctx)

	<-ctx.Done()
}
