package usecase

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

type UseCase struct {
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

func New(storage storage.Storage, logger Logger) *UseCase {
	return &UseCase{
		storage: storage,
		logg:    logger,
	}
}

func (a *UseCase) CreateEvent(ctx context.Context, id, title string) error {
	return nil
}
