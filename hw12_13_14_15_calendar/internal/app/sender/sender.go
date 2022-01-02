package appsender

import (
	"context"
	"fmt"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
)

const ThreadsCount = 2

type Sender struct {
	consumer Consumer
	logger   app.Logger
}

func New(consumer Consumer, logger app.Logger) *Sender {
	return &Sender{
		consumer: consumer,
		logger:   logger,
	}
}

func (s *Sender) Send(ctx context.Context) {
	s.consumer.Consume(ctx, func(ctx context.Context, ch <-chan Event) {
		select {
		case event := <-ch:
			s.logger.Info(fmt.Sprintf(
				"Send notification event '%s', id '%s', for user '%s'",
				event.Title,
				event.ID.String(),
				event.UserID.String(),
			))
		case <-ctx.Done():
			return
		}
	}, ThreadsCount)
}
