package notifier

import (
	"context"
	"log"
	"time"

	appscheduler "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

// TODO add logger.
type EventsNotifier struct {
	storage  storage.Storage
	producer appscheduler.Producer
}

func New(storage storage.Storage, producer appscheduler.Producer) *EventsNotifier {
	return &EventsNotifier{
		storage:  storage,
		producer: producer,
	}
}

func (n *EventsNotifier) NotifyEvents(ctx context.Context, duration time.Duration) error {
	// TODO change ticker duration
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			events, err := n.storage.GetUpcomingEvents(ctx, 5*time.Hour)
			if err != nil {
				log.Println(err)
			}

			for key := range events {
				if err := n.producer.Produce(ctx, &events[key]); err != nil {
					log.Println(err)
					continue
				}

				events[key].NotifiedAt = time.Now()

				if err := n.storage.Update(ctx, &events[key]); err != nil {
					log.Println(err)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
