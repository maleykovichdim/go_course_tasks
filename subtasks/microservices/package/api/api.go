package api

import "context"

type Dispatcher interface {
	Write(ctx context.Context, m string) error
	ReadConfirmed(ctx context.Context) (string, error)
}
