package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar/usecase"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/logger"
	internalgrpc "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/server/grpc"
	internalhttp "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/server/http"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar.dev.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", config.Logger)

	logger, err := logger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.New(config.Storage)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	err = storage.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close(ctx)

	calendar := usecase.New(storage, logger)
	httpserver := internalhttp.NewServer(config.Server, calendar, logger)
	grpcserver := internalgrpc.NewServer(config.Grpc, calendar, logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpserver.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	go func() {
		defer cancel()
		if err := httpserver.Start(ctx); err != nil {
			logger.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		defer cancel()
		if err := grpcserver.Start(ctx); err != nil {
			logger.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	<-ctx.Done()
}
