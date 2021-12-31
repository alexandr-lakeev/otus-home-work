package appscheduler

import (
	"context"
	"time"
)

type EventsNotifier interface {
	NotifyEvents(ctx context.Context, duration time.Duration) error
}
