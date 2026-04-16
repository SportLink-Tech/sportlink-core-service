package slices

// Map applies a transformation function to each element of a slice and returns a new slice with the results.
func Map[T any, E any](slice []T, fn func(T) E) []E {
	result := make([]E, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Filter returns a new slice containing only the elements for which the predicate function returns true.
func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Contains reports whether target is present in slice.
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
