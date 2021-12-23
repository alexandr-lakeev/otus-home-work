package usecase

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
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
	event := models.Event{
		ID:          command.ID,
		UserID:      command.UserID,
		Title:       command.Title,
		Date:        command.Date,
		Duration:    command.Duration,
		Description: command.Description,
	}

	if err := u.checkDateBusyForEvent(ctx, &event); err != nil {
		return err
	}

	return u.storage.Add(ctx, &event)
}

func (u *UseCase) checkDateBusyForEvent(ctx context.Context, event *models.Event) error {
	events, err := u.storage.GetList(ctx, event.UserID, event.Date, event.Date.Add(event.Duration))
	if err != nil {
		return err
	}

	if len(events) != 0 {
		for _, e := range events {
			if event.ID != e.ID {
				return domain.ErrDateBusy
			}
		}
	}

	return nil
}
