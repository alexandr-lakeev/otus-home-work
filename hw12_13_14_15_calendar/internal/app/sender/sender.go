package appsender

import (
	"context"
	"fmt"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
)

const ThreadsCount = 2

type Sender struct {
	consumer Consumer
	producer app.Producer
	logger   app.Logger
}

func New(consumer Consumer, producer app.Producer, logger app.Logger) *Sender {
	return &Sender{
		consumer: consumer,
		producer: producer,
		logger:   logger,
	}
}

func (s *Sender) Send(ctx context.Context) {
	s.consumer.Consume(ctx, func(ctx context.Context, ch <-chan Event) {
		for {
			select {
			case event := <-ch:
				s.logger.Info(fmt.Sprintf(
					"Send notification event '%s', id '%s', for user '%s'",
					event.Title,
					event.ID.String(),
					event.UserID.String(),
				))

				if err := s.producer.Produce(ctx, event); err != nil {
					s.logger.Error("Error on produce event notified: " + err.Error())
				}
			case <-ctx.Done():
				s.logger.Info("Stop sender...")
				return
			}
		}
	}, ThreadsCount)
}
