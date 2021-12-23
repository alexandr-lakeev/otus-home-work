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
	t.Run("create event", func(t *testing.T) {
		logger, _ := logger.New(config.LoggerConf{
			Level: "INFO",
		})

		usecase := New(memorystorage.New(), logger)

		err := usecase.CreateEvent(context.Background(), &app.CreateEventCommand{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Title:  "Dummy event",
			Date:   time.Now(),
		})

		require.NoError(t, err)
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

		var errDateBusy = domain.ErrDateBusy

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
