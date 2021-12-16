package storage

import (
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/google/uuid"
)

type Storage interface {
	Get(id uuid.UUID) (*models.Event, error)
	Save(event *models.Event) error
	GetList(from, to time.Time) ([]models.Event, error)
}
