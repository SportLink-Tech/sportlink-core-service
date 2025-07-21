package slice

// Contains returns true if the given slice contains the item `target`,
// using a custom comparator function to determine equality.
func Contains[T any](slice []T, target T, comparator func(a, b T) bool) bool {
	for _, elem := range slice {
		if comparator(elem, target) {
			return true
		}
	}
	return false
}
