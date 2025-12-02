package matchannouncement

import (
	"fmt"
	"sportlink/api/domain/common"
)

// CategoryRange represents the admitted categories for a match
// Can be a specific category, a range, or multiple categories
type CategoryRange struct {
	Type       RangeType         // Range type
	Categories []common.Category // Specific list of categories (e.g: L5, L7)
	MinLevel   common.Category   // Minimum level for ranges (e.g: >= L5)
	MaxLevel   common.Category   // Maximum level for ranges (e.g: <= L5)
}

type RangeType string

const (
	RangeTypeSpecific    RangeType = "SPECIFIC"     // Specific categories (e.g: only L5 and L7)
	RangeTypeGreaterThan RangeType = "GREATER_THAN" // Greater than or equal (e.g: >= L5)
	RangeTypeLessThan    RangeType = "LESS_THAN"    // Less than or equal (e.g: <= L5)
	RangeTypeBetween     RangeType = "BETWEEN"      // Between two levels (e.g: L3 to L6)
)

// NewSpecificCategories creates a range with specific categories
func NewSpecificCategories(categories []common.Category) CategoryRange {
	return CategoryRange{
		Type:       RangeTypeSpecific,
		Categories: categories,
	}
}

// NewGreaterThanCategory creates a range of categories >= min
func NewGreaterThanCategory(min common.Category) CategoryRange {
	return CategoryRange{
		Type:     RangeTypeGreaterThan,
		MinLevel: min,
	}
}

// NewLessThanCategory creates a range of categories <= max
func NewLessThanCategory(max common.Category) CategoryRange {
	return CategoryRange{
		Type:     RangeTypeLessThan,
		MaxLevel: max,
	}
}

// NewBetweenCategories creates a range of categories between min and max
func NewBetweenCategories(min, max common.Category) (CategoryRange, error) {
	if min > max {
		return CategoryRange{}, fmt.Errorf("min category cannot be greater than max category")
	}
	return CategoryRange{
		Type:     RangeTypeBetween,
		MinLevel: min,
		MaxLevel: max,
	}, nil
}

// Admits checks if a category is admitted in this range
func (cr CategoryRange) Admits(category common.Category) bool {
	switch cr.Type {
	case RangeTypeSpecific:
		for _, c := range cr.Categories {
			if c == category {
				return true
			}
		}
		return false
	case RangeTypeGreaterThan:
		return category >= cr.MinLevel
	case RangeTypeLessThan:
		return category <= cr.MaxLevel
	case RangeTypeBetween:
		return category >= cr.MinLevel && category <= cr.MaxLevel
	default:
		return false
	}
}
