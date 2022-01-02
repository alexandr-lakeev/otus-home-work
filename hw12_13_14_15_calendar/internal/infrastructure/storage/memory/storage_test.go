package memorystorage

import (
	"context"
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

		ctx := context.Background()
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

		require.NoError(t, s.Add(ctx, event))

		eventFromStorage, err := s.Get(ctx, eventID)

		require.NoError(t, err)
		require.Equal(t, eventFromStorage, event)
	})

	t.Run("not found", func(t *testing.T) {
		s := New()

		ctx := context.Background()
		eventID := uuid.New()

		event := &models.Event{
			ID: eventID,
		}

		require.NoError(t, s.Add(ctx, event))

		_, err := s.Get(ctx, uuid.New())

		require.ErrorIs(t, domain.ErrEventNotFound, err)
	})

	t.Run("list", func(t *testing.T) {
		s := New()

		ctx := context.Background()
		userID := uuid.New()

		date1, _ := time.Parse("2006-01-02 15:04:05", "2021-12-18 12:00:00")
		event1 := &models.Event{
			ID:          uuid.New(),
			Title:       "New Event 1",
			Date:        date1,
			Duration:    2 * time.Hour,
			Description: "Some awesome event 1",
			UserID:      userID,
		}
		require.NoError(t, s.Add(ctx, event1))

		date2, _ := time.Parse("2006-01-02 15:04:05", "2021-12-19 13:00:00")
		event2 := &models.Event{
			ID:          uuid.New(),
			Title:       "New Event 2",
			Date:        date2,
			Duration:    2 * time.Hour,
			Description: "Some awesome event 2",
			UserID:      userID,
		}
		require.NoError(t, s.Add(ctx, event2))

		date3, _ := time.Parse("2006-01-02 15:04:05", "2021-12-20 14:00:00")
		event3 := &models.Event{
			ID:          uuid.New(),
			Title:       "New Event 3",
			Date:        date3,
			Duration:    2 * time.Hour,
			Description: "Some awesome event 3",
			UserID:      userID,
		}
		require.NoError(t, s.Add(ctx, event3))

		// Another user event
		date4, _ := time.Parse("2006-01-02 15:04:05", "2021-12-18 12:00:00")
		event4 := &models.Event{
			ID:          uuid.New(),
			Title:       "New Event 2",
			Date:        date4,
			Duration:    2 * time.Hour,
			Description: "Some awesome event",
			UserID:      uuid.New(),
		}
		require.NoError(t, s.Add(ctx, event4))

		from, _ := time.Parse("2006-01-02 15:04:05", "2021-12-18 00:00:00")
		to, _ := time.Parse("2006-01-02 15:04:05", "2021-12-20 00:00:00")
		list, err := s.GetList(ctx, userID, from, to)

		require.NoError(t, err)
		require.Len(t, list, 2)
	})
}
