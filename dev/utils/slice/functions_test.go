package slice_test

import (
	"sportlink/dev/utils/slice"
	"testing"
)

type person struct {
	ID   int
	Name string
}

func TestContains(t *testing.T) {
	tests := []struct {
		name       string
		slice      []person
		target     person
		comparator func(a, b person) bool
		expected   bool
	}{
		{
			name: "found by ID",
			slice: []person{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			target:   person{ID: 2},
			expected: true,
			comparator: func(a, b person) bool {
				return a.ID == b.ID
			},
		},
		{
			name: "not found",
			slice: []person{
				{ID: 1, Name: "Alice"},
			},
			target:   person{ID: 3},
			expected: false,
			comparator: func(a, b person) bool {
				return a.ID == b.ID
			},
		},
		{
			name: "found by name",
			slice: []person{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			target:   person{Name: "Alice"},
			expected: true,
			comparator: func(a, b person) bool {
				return a.Name == b.Name
			},
		},
		{
			name:     "empty slice",
			slice:    []person{},
			target:   person{ID: 1},
			expected: false,
			comparator: func(a, b person) bool {
				return a.ID == b.ID
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := slice.Contains(tc.slice, tc.target, tc.comparator)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
