package parser

import (
	"fmt"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	"strconv"
	"strings"
	"time"
)

// QueryParser defines the interface for parsing query parameters
type QueryParser interface {
	ParseSports(sportsQuery string) ([]common.Sport, error)
	ParseCategories(categoriesQuery string) ([]common.Category, error)
	ParseStatuses(statusesQuery string) ([]matchannouncement.Status, error)
	ParseDate(dateQuery string) (time.Time, error)
	ParseLocation(country, province, locality string) *matchannouncement.Location
}

// DefaultQueryParser implements QueryParser interface
type DefaultQueryParser struct{}

// NewQueryParser creates a new instance of DefaultQueryParser
func NewQueryParser() QueryParser {
	return &DefaultQueryParser{}
}

// ParseSports parses a comma-separated string of sports into a slice of Sport
func (p *DefaultQueryParser) ParseSports(sportsQuery string) ([]common.Sport, error) {
	if sportsQuery == "" {
		return nil, nil
	}

	sportStrings := strings.Split(sportsQuery, ",")
	sports := make([]common.Sport, 0, len(sportStrings))

	for _, sportStr := range sportStrings {
		trimmed := strings.TrimSpace(sportStr)
		if trimmed == "" {
			continue
		}
		sports = append(sports, common.Sport(trimmed))
	}

	return sports, nil
}

// ParseCategories parses a comma-separated string of category numbers into a slice of Category
func (p *DefaultQueryParser) ParseCategories(categoriesQuery string) ([]common.Category, error) {
	if categoriesQuery == "" {
		return nil, nil
	}

	categoryStrings := strings.Split(categoriesQuery, ",")
	categories := make([]common.Category, 0, len(categoryStrings))

	for _, catStr := range categoryStrings {
		trimmed := strings.TrimSpace(catStr)
		if trimmed == "" {
			continue
		}

		catInt, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid category format: %s", trimmed)
		}

		category, err := common.GetCategory(catInt)
		if err != nil {
			return nil, fmt.Errorf("invalid category value: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// ParseStatuses parses a comma-separated string of statuses into a slice of Status
func (p *DefaultQueryParser) ParseStatuses(statusesQuery string) ([]matchannouncement.Status, error) {
	if statusesQuery == "" {
		return nil, nil
	}

	statusStrings := strings.Split(statusesQuery, ",")
	statuses := make([]matchannouncement.Status, 0, len(statusStrings))

	for _, statusStr := range statusStrings {
		trimmed := strings.TrimSpace(statusStr)
		if trimmed == "" {
			continue
		}

		status, err := matchannouncement.ParseStatus(trimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid status value: %w", err)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// ParseDate parses a date string in YYYY-MM-DD format into a time.Time
func (p *DefaultQueryParser) ParseDate(dateQuery string) (time.Time, error) {
	if dateQuery == "" {
		return time.Time{}, nil
	}

	date, err := time.Parse("2006-01-02", dateQuery)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	return date, nil
}

// ParseLocation creates a Location from country, province, and locality strings
func (p *DefaultQueryParser) ParseLocation(country, province, locality string) *matchannouncement.Location {
	if country == "" && province == "" && locality == "" {
		return nil
	}

	return &matchannouncement.Location{
		Country:  country,
		Province: province,
		Locality: locality,
	}
}
