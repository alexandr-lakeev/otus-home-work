package memorystorage

import (
	"sync"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]*models.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]*models.Event),
	}
}

func (s *Storage) Get(id uuid.UUID) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[id]
	if !ok {
		return nil, storage.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) Save(event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return nil
}

func (s *Storage) GetList(from, to time.Time) ([]models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return []models.Event{}, nil
}