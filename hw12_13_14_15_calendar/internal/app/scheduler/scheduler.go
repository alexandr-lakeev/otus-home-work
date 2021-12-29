package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
)

type UseCase struct {
	storage storage.Storage
}

func New(storage storage.Storage) *UseCase {
	return &UseCase{
		storage: storage,
	}
}

func (u *UseCase) Notify(ctx context.Context, duration time.Duration) error {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick")
		case <-ctx.Done():
			return nil
		}
	}
}

func (u *UseCase) ClearOldEvent(ctx context.Context, duration time.Duration) error {
	return nil
}
