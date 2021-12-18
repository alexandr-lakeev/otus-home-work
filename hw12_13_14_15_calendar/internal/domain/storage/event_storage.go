package storage

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
)

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Get(ctx context.Context, id models.ID) (*models.Event, error)
	Add(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.Event) error
	GetList(ctx context.Context, from, to time.Time) ([]models.Event, error)
}
