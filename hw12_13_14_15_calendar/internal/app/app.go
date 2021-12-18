package app

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

type App struct {
	logg    Logger
	storage storage.Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Panic(msg string)
}

func New(storage storage.Storage, logger Logger) *App {
	return &App{
		storage: storage,
		logg:    logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return nil
}

// TODO
