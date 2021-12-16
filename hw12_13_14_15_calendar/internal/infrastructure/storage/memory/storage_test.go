package memorystorage

import (
	"testing"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("save and get", func(t *testing.T) {
		s := New()

		eventID := uuid.New()
		userID := uuid.New()

		event := &models.Event{
			ID:          eventID,
			Title:       "New Event",
			Date:        time.Now(),
			Duration:    2 * time.Hour,
			Description: "Some awesome event",
			UserID:      userID,
		}

		require.NoError(t, s.Save(event))

		eventFromStorage, err := s.Get(eventID)

		require.NoError(t, err)
		require.Equal(t, eventFromStorage, event)
	})

	t.Run("not found", func(t *testing.T) {
		s := New()
		eventID := uuid.New()

		event := &models.Event{
			ID: eventID,
		}

		require.NoError(t, s.Save(event))

		_, err := s.Get(uuid.New())

		require.ErrorIs(t, domain.ErrEventNotFound, err)
	})

	t.Run("date busy", func(t *testing.T) {
		s := New()

		userID := uuid.New()
		date := time.Now()

		event1 := &models.Event{
			ID:          uuid.New(),
			Title:       "New Event 1",
			Date:        date,
			Duration:    2 * time.Hour,
			Description: "Some awesome event",
			UserID:      userID,
		}

		event2 := &models.Event{
			ID:          uuid.New(),
			Title:       "New Event 2",
			Date:        date,
			Duration:    2 * time.Hour,
			Description: "Some awesome event",
			UserID:      userID,
		}

		require.NoError(t, s.Save(event1))
		require.ErrorIs(t, domain.ErrDateBusy, s.Save(event2))
	})
}
