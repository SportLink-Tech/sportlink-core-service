package common

//go:generate stringer -type=Category
type Category int

const (
	L1 Category = iota + 1
	L2
	L3
	L4
	L5
	L6
	L7
)