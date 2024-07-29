package application

// UseCase is a generic interface that defines a single method Invoke,
// taking an input of type T and returning an output of type O along with an error.go.
type UseCase[T any, O any] interface {
	Invoke(input T) (*O, error)
}
