package app

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
)

type UseCase interface {
	CreateEvent(context.Context, *CreateEventCommand) error
}

type CreateEventCommand struct {
	ID          models.ID
	UserID      models.ID
	Title       string
	Description string
	Date        time.Time
	Duration    time.Duration
}
