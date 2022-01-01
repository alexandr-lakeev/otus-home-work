package cleaner

import (
	"context"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

type EventsCleaner struct {
	storage storage.Storage
	logger  app.Logger
}

func New(storage storage.Storage, logger app.Logger) *EventsCleaner {
	return &EventsCleaner{
		storage: storage,
		logger:  logger,
	}
}

func (c *EventsCleaner) ClearEvents(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			c.logger.Info("Start old events cleaning")

			if err := c.storage.DeleteEvents(ctx, duration); err != nil {
				c.logger.Error(err.Error())
			} else {
				c.logger.Info("Old events cleaned")
			}
		case <-ctx.Done():
			return
		}
	}
}
