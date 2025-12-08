package application

import "context"

// UseCase is a generic interface that defines a single method Invoke,
// taking a context and an input of type T and returning an output of type O along with an error.
type UseCase[T any, O any] interface {
	Invoke(ctx context.Context, input T) (*O, error)
}
