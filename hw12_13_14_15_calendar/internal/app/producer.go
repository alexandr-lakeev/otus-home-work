package app

import "context"

type Producer interface {
	Produce(ctx context.Context, data interface{}) error
}
