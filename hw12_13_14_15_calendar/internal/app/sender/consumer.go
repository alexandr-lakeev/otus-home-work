package appsender

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
)

type Consumer interface {
	Consume(ctx context.Context, fn Worker, threads int) error
}

type Worker func(context.Context, <-chan Event)

type Event struct {
	ID     models.ID
	UserID models.ID
	Title  string
}
