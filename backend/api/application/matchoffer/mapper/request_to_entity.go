package mapper

import (
	"fmt"
	"sportlink/api/application/matchoffer/request"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchoffer"
	"time"
)

// parseDateTime tries to parse datetime with multiple formats
func parseDateTime(dateTimeStr string, loc *time.Location) (time.Time, error) {
	// Try RFC3339 first (with timezone)
	t, err := time.Parse(time.RFC3339, dateTimeStr)
	if err == nil {
		return t.In(loc), nil
	}

	// Try without timezone (format from frontend: 2025-12-05T13:00:00)
	t, err = time.ParseInLocation("2006-01-02T15:04:05", dateTimeStr, loc)
	if err == nil {
		return t, nil
	}

	// Try with seconds (alternative format)
	t, err = time.ParseInLocation("2006-01-02T15:04", dateTimeStr, loc)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateTimeStr)
}

func CreationRequestToEntity(req request.NewMatchOfferRequest, ownerAccountID string) (matchoffer.Entity, error) {
	// Build location first to get timezone
	var location matchoffer.Location
	if req.Location.Latitude != nil && req.Location.Longitude != nil {
		location = matchoffer.NewLocationWithCoords(req.Location.Country, req.Location.Province, req.Location.Locality, *req.Location.Latitude, *req.Location.Longitude)
	} else {
		location = matchoffer.NewLocation(req.Location.Country, req.Location.Province, req.Location.Locality)
	}
	tz := location.GetTimezone()

	// Parse day in the location's timezone (not UTC) to avoid off-by-one day errors
	day, err := time.ParseInLocation("2006-01-02", req.Day, tz)
	if err != nil {
		return matchoffer.Entity{}, fmt.Errorf("invalid day format: %w", err)
	}

	// Parse time slot (try with and without timezone)
	startTime, err := parseDateTime(req.TimeSlot.StartTime, tz)
	if err != nil {
		return matchoffer.Entity{}, fmt.Errorf("invalid start time format: %w", err)
	}

	endTime, err := parseDateTime(req.TimeSlot.EndTime, tz)
	if err != nil {
		return matchoffer.Entity{}, fmt.Errorf("invalid end time format: %w", err)
	}

	timeSlot, err := matchoffer.NewTimeSlot(startTime, endTime)
	if err != nil {
		return matchoffer.Entity{}, err
	}

	// Build category range
	var categoryRange matchoffer.CategoryRange
	switch req.AdmittedCategories.Type {
	case "SPECIFIC":
		categories := make([]common.Category, len(req.AdmittedCategories.Categories))
		for i, c := range req.AdmittedCategories.Categories {
			cat, err := common.GetCategory(c)
			if err != nil {
				return matchoffer.Entity{}, err
			}
			categories[i] = cat
		}
		categoryRange = matchoffer.NewSpecificCategories(categories)
	case "GREATER_THAN":
		minCat, err := common.GetCategory(req.AdmittedCategories.MinLevel)
		if err != nil {
			return matchoffer.Entity{}, err
		}
		categoryRange = matchoffer.NewGreaterThanCategory(minCat)
	case "LESS_THAN":
		maxCat, err := common.GetCategory(req.AdmittedCategories.MaxLevel)
		if err != nil {
			return matchoffer.Entity{}, err
		}
		categoryRange = matchoffer.NewLessThanCategory(maxCat)
	case "BETWEEN":
		minCat, err := common.GetCategory(req.AdmittedCategories.MinLevel)
		if err != nil {
			return matchoffer.Entity{}, err
		}
		maxCat, err := common.GetCategory(req.AdmittedCategories.MaxLevel)
		if err != nil {
			return matchoffer.Entity{}, err
		}
		categoryRange, err = matchoffer.NewBetweenCategories(minCat, maxCat)
		if err != nil {
			return matchoffer.Entity{}, err
		}
	default:
		return matchoffer.Entity{}, fmt.Errorf("invalid category range type: %s", req.AdmittedCategories.Type)
	}

	// Build entity
	sport := common.Sport(req.Sport)
	dayInTz := day.In(tz)
	createdAt := time.Now().In(tz)

	return matchoffer.NewMatchOffer(
		req.TeamName,
		sport,
		dayInTz,
		timeSlot,
		location,
		categoryRange,
		matchoffer.StatusPending,
		createdAt,
		ownerAccountID,
		req.Capacity,
	), nil
}
