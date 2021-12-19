package usecase

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

type UseCase struct {
	logg    app.Logger
	storage storage.Storage
}

func New(storage storage.Storage, logger app.Logger) *UseCase {
	return &UseCase{
		storage: storage,
		logg:    logger,
	}
}

func (a *UseCase) CreateEvent(ctx context.Context, id, title string) error {
	return nil
}
