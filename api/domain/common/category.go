package common

import "fmt"

//go:generate stringer -type=Category
type Category int

const (
	Unranked Category = 0
	L1       Category = iota + 1
	L2
	L3
	L4
	L5
	L6
	L7
	MaxCategory
)

func GetCategory(value int) (Category, error) {
	if value < int(L1) || value >= int(MaxCategory) {
		return 0, fmt.Errorf("invalid category value: %d", value)
	}
	return Category(value), nil
}
