package memorystorage

import (
	"testing"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("save and get", func(t *testing.T) {
		s := New()

		eventId := uuid.New()
		userId := uuid.New()

		event := &models.Event{
			ID:          eventId,
			Title:       "New Event",
			Date:        time.Now(),
			Duration:    2 * time.Hour,
			Description: "Some awesome event",
			UserId:      userId,
		}

		require.NoError(t, s.Save(event))

		eventFromStorage, err := s.Get(eventId)

		require.NoError(t, err)
		require.Equal(t, eventFromStorage, event)
	})

	t.Run("not found", func(t *testing.T) {
		s := New()
		eventId := uuid.New()

		event := &models.Event{
			ID: eventId,
		}

		require.NoError(t, s.Save(event))

		_, err := s.Get(uuid.New())

		require.ErrorIs(t, storage.ErrEventNotFound, err)
	})
}
