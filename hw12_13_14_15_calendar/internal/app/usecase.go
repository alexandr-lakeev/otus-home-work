package app

import "context"

type UseCase interface {
	CreateEvent(ctx context.Context, id, title string) error
}
