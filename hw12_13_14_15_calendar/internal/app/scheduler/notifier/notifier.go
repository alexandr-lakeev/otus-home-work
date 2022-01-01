package notifier

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	appscheduler "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

// TODO add logger.
type EventsNotifier struct {
	storage  storage.Storage
	producer appscheduler.Producer
	logger   app.Logger
}

func New(storage storage.Storage, producer appscheduler.Producer, logger app.Logger) *EventsNotifier {
	return &EventsNotifier{
		storage:  storage,
		producer: producer,
		logger:   logger,
	}
}

func (n *EventsNotifier) NotifyEvents(ctx context.Context, duration time.Duration) {
	// TODO change ticker duration
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			n.logger.Info("Start events notification")

			events, err := n.storage.GetUpcomingEvents(ctx, 5*time.Hour)
			if err != nil {
				n.logger.Error(err.Error())
			}

			for key := range events {
				events[key].NotifiedAt = time.Now()

				if err := n.producer.Produce(ctx, &events[key]); err != nil {
					n.logger.Error(err.Error())
					continue
				}

				if err := n.storage.Update(ctx, &events[key]); err != nil {
					n.logger.Error(err.Error())
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
