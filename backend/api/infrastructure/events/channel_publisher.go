package events

import "context"

// ChannelPublisher is a generic Publisher backed by a buffered Go channel.
// It is safe to use from multiple goroutines.
type ChannelPublisher[T any] struct {
	ch chan T
}

func NewChannelPublisher[T any](bufferSize int) *ChannelPublisher[T] {
	return &ChannelPublisher[T]{ch: make(chan T, bufferSize)}
}

// Ch returns the read-only event channel for consumers.
func (p *ChannelPublisher[T]) Ch() <-chan T {
	return p.ch
}

func (p *ChannelPublisher[T]) Publish(_ context.Context, event T) error {
	p.ch <- event
	return nil
}
