// Code generated by "stringer -type=Category"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[L1-1]
	_ = x[L2-2]
	_ = x[L3-3]
	_ = x[L4-4]
	_ = x[L5-5]
	_ = x[L6-6]
	_ = x[L7-7]
}

const _Category_name = "L1L2L3L4L5L6L7"

var _Category_index = [...]uint8{0, 2, 4, 6, 8, 10, 12, 14}

func (i Category) String() string {
	i -= 1
	if i < 0 || i >= Category(len(_Category_index)-1) {
		return "Category(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Category_name[_Category_index[i]:_Category_index[i+1]]
}