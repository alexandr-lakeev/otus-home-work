package appscheduler

import (
	"context"
	"time"
)

type EventsCleaner interface {
	ClearEvents(ctx context.Context, duration time.Duration) error
}
