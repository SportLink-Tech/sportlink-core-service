package matchoffer

import (
	"fmt"
	"sportlink/api/domain/common"
	"strings"
)

var countryAbbreviations = map[string]string{
	"Argentina": "ARG",
	"Brasil":    "BRA",
	"Brazil":    "BRA",
	"Uruguay":   "URU",
	"Chile":     "CHI",
	"Colombia":  "COL",
	"España":    "ESP",
	"Spain":     "ESP",
}

// BuildTitle constructs a human-readable title from the given sport, categories and location.
// Format: "<Sport> · <Categories> · <Locality>, <CountryCode>"
// Example: "Paddle · L3-L5 · Palermo, ARG"
func BuildTitle(sport common.Sport, categories CategoryRange, location Location) string {
	parts := []string{string(sport)}

	if cat := formatCategories(categories); cat != "" {
		parts = append(parts, cat)
	}

	country := location.Country
	if abbr, ok := countryAbbreviations[country]; ok {
		country = abbr
	}
	parts = append(parts, fmt.Sprintf("%s, %s", location.Locality, country))

	return strings.Join(parts, " · ")
}

func formatCategories(cr CategoryRange) string {
	switch cr.Type {
	case RangeTypeSpecific:
		if len(cr.Categories) == 0 {
			return ""
		}
		labels := make([]string, len(cr.Categories))
		for i, c := range cr.Categories {
			labels[i] = fmt.Sprintf("L%d", c)
		}
		return strings.Join(labels, ", ")
	case RangeTypeGreaterThan:
		return fmt.Sprintf("L%d+", cr.MinLevel)
	case RangeTypeLessThan:
		return fmt.Sprintf("L%d-", cr.MaxLevel)
	case RangeTypeBetween:
		return fmt.Sprintf("L%d-L%d", cr.MinLevel, cr.MaxLevel)
	default:
		return ""
	}
}
