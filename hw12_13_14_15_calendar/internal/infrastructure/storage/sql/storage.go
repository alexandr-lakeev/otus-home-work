package sqlstorage

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
)

type Storage struct {
	// conn sql.Conn
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (*Storage) Get(id string) (*models.Event, error) {
	return &models.Event{}, nil
}

func (*Storage) Save(event *models.Event) error {
	return nil
}

func (*Storage) GetList(from, to time.Time) ([]models.Event, error) {
	return []models.Event{}, nil
}
