package notifier

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	appscheduler "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

const RemindTickerDuration = 1 * time.Minute

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
	ticker := time.NewTicker(RemindTickerDuration)

	for {
		select {
		case <-ticker.C:
			n.logger.Info("Start events notification")

			events, err := n.storage.GetUpcomingEvents(ctx, 5*time.Hour)
			if err != nil {
				n.logger.Error(err.Error())
			}

			for key, event := range events {
				n.logger.Info("Event notification " + event.ID.String())

				events[key].NotifiedAt = time.Now()

				if err := n.producer.Produce(ctx, &events[key]); err != nil {
					n.logger.Error(err.Error())
					continue
				}

				if err := n.storage.Update(ctx, &events[key]); err != nil {
					n.logger.Error(err.Error())
				}
			}

			n.logger.Info("Complete events notification")
		case <-ctx.Done():
			return
		}
	}
}
