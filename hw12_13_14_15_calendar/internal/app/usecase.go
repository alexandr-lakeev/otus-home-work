package app

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
)

type UseCase interface {
	CreateEvent(context.Context, *CreateEventCommand) error
	GetEvent(context.Context, *GetEventQuery) (*models.Event, error)
	GetList(context.Context, *GetListQuery) ([]models.Event, error)
	UpdateEvent(context.Context, *UpdateEventCommand) error
}

type CreateEventCommand struct {
	ID          models.ID
	UserID      models.ID
	Title       string
	Description string
	Date        time.Time
	Duration    time.Duration
}

type GetEventQuery struct {
	ID     models.ID
	UserID models.ID
}

type GetListQuery struct {
	UserID models.ID
	From   time.Time
	To     time.Time
}

type UpdateEventCommand struct {
	ID          models.ID
	UserID      models.ID
	Title       string
	Description string
	Date        time.Time
	Duration    time.Duration
}
