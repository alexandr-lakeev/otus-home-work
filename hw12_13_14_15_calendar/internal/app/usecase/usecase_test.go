package usecase

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/logger"
	memorystorage "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUseCase(t *testing.T) {
	t.Run("create and get event", func(t *testing.T) {
		logger, _ := logger.New(config.LoggerConf{
			Level: "INFO",
		})

		id := uuid.New()
		userID := uuid.New()
		title := "Dummy event"
		date := time.Now()
		usecase := New(memorystorage.New(), logger)

		err := usecase.CreateEvent(context.Background(), &app.CreateEventCommand{
			ID:     id,
			UserID: userID,
			Title:  title,
			Date:   date,
		})

		require.NoError(t, err)

		event, err := usecase.GetEvent(context.Background(), &app.GetEventQuery{
			ID:     id,
			UserID: userID,
		})

		require.NoError(t, err)
		require.Equal(t, id, event.ID)
		require.Equal(t, title, event.Title)
		require.Equal(t, date, event.Date)
	})

	t.Run("date busy", func(t *testing.T) {
		logger, _ := logger.New(config.LoggerConf{
			Level: "INFO",
		})

		usecase := New(memorystorage.New(), logger)

		userID := uuid.New()
		beginDate := time.Now()

		err := usecase.CreateEvent(context.Background(), &app.CreateEventCommand{
			ID:       uuid.New(),
			UserID:   userID,
			Title:    "Dummy event",
			Date:     beginDate,
			Duration: 10 * time.Minute,
		})

		require.NoError(t, err)

		err = usecase.CreateEvent(context.Background(), &app.CreateEventCommand{
			ID:       uuid.New(),
			UserID:   userID,
			Title:    "Dummy event",
			Date:     beginDate.Add(5 * time.Minute),
			Duration: 5 * time.Minute,
		})

		errDateBusy := domain.ErrDateBusy

		require.Error(t, err)
		require.ErrorIs(t, err, errDateBusy)

		err = usecase.CreateEvent(context.Background(), &app.CreateEventCommand{
			ID:       uuid.New(),
			UserID:   userID,
			Title:    "Dummy event",
			Date:     beginDate.Add(-5 * time.Minute),
			Duration: 10 * time.Minute,
		})

		log.Println("Date", beginDate.Add(-5*time.Minute))

		require.Error(t, err)
		require.ErrorIs(t, err, errDateBusy)
	})
}
