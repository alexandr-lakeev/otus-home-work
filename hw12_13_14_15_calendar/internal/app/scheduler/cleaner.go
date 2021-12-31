package appscheduler

import (
	"context"
	"time"
)

type EventsCleaner interface {
	ClearEvent(ctx context.Context, duration time.Duration) error
}
