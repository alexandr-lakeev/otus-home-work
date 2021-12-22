package usecase

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
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

func (u *UseCase) CreateEvent(ctx context.Context, command *app.CreateEventCommand) error {
	// TODO check date for busy
	return u.storage.Add(ctx, &models.Event{
		ID:          command.ID,
		UserID:      command.UserID,
		Title:       command.Title,
		Date:        command.Date,
		Duration:    command.Duration,
		Description: command.Description,
	})
}
