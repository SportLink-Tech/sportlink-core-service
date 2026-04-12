package events

import "context"

// Publisher is a generic interface for publishing domain events.
// T is the event payload type.
type Publisher[T any] interface {
	Publish(ctx context.Context, event T) error
}
