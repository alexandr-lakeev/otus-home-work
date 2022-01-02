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

func (u *UseCase) GetEvent(ctx context.Context, query *app.GetEventQuery) (*models.Event, error) {
	event, err := u.storage.Get(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	if event.UserID != query.UserID {
		return nil, domain.ErrPremissionDenied
	}

	return event, nil
}

func (u *UseCase) GetList(ctx context.Context, query *app.GetListQuery) ([]models.Event, error) {
	return u.storage.GetList(ctx, query.UserID, query.From, query.To)
}

func (u *UseCase) UpdateEvent(ctx context.Context, command *app.UpdateEventCommand) error {
	event, err := u.storage.Get(ctx, command.ID)
	if err != nil {
		return err
	}

	if event.UserID != command.UserID {
		return domain.ErrPremissionDenied
	}

	event.Title = command.Title
	event.Date = command.Date
	event.Duration = command.Duration
	event.Description = command.Description

	if err := u.checkDateBusyForEvent(ctx, event); err != nil {
		return err
	}

	return u.storage.Update(ctx, event)
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
