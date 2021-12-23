package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
)

type Storage struct {
	mu     sync.RWMutex
	events map[models.ID]*models.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[models.ID]*models.Event),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) Get(ctx context.Context, id models.ID) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[id]
	if !ok {
		return nil, domain.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) Add(ctx context.Context, event *models.Event) error {
	return s.Update(ctx, event)
}

func (s *Storage) Update(ctx context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return nil
}

func (s *Storage) GetList(ctx context.Context, userID models.ID, from, to time.Time) ([]models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Event
	for _, e := range s.events {
		if e.UserID != userID {
			continue
		}

		from2 := e.Date
		to2 := e.Date.Add(e.Duration)

		if (from.Equal(to2) || from.Before(to2)) && (to.Equal(from2) || to.After(from2)) {
			result = append(result, *e)
		}
	}

	return result, nil
}
