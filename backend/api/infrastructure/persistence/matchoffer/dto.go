package matchoffer

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchoffer"
	"time"
)

type Dto struct {
	EntityId       string   `dynamodbav:"EntityId"`                 // "Entity#MatchOffer"
	Id             string   `dynamodbav:"Id"`                       // Generated UUID
	TeamName       string   `dynamodbav:"TeamName"`                 // Team name
	Sport          string   `dynamodbav:"Sport"`                    // Sport type
	Day            int64    `dynamodbav:"Day"`                      // Unix timestamp of the day
	StartTime      int64    `dynamodbav:"StartTime"`                // Unix timestamp of start time
	EndTime        int64    `dynamodbav:"EndTime"`                  // Unix timestamp of end time
	Country        string   `dynamodbav:"Country"`                  // Location country
	Province       string   `dynamodbav:"Province"`                 // Location province
	Locality       string   `dynamodbav:"Locality"`                 // Location locality
	GeohashPrefix  *string  `dynamodbav:"GeohashPrefix,omitempty"`  // Geohash prefix (precision 3) for GSI — absent when no coords
	Latitude       *float64 `dynamodbav:"Latitude,omitempty"`       // GPS latitude — absent when no coords
	Longitude      *float64 `dynamodbav:"Longitude,omitempty"`      // GPS longitude — absent when no coords
	RangeType      string   `dynamodbav:"RangeType"`                // Category range type
	Categories     []int    `dynamodbav:"Categories"`               // Specific categories
	MinLevel       int      `dynamodbav:"MinLevel"`                 // Minimum category level
	MaxLevel       int      `dynamodbav:"MaxLevel"`                 // Maximum category level
	Status         string   `dynamodbav:"Status"`                   // Offer status
	CreatedAt      int64    `dynamodbav:"CreatedAt"`                // Unix timestamp of creation
	ExpiresAt      int64    `dynamodbav:"ExpiresAt"`                // TTL for DynamoDB (Unix timestamp)
	OwnerAccountId string   `dynamodbav:"OwnerAccountId,omitempty"` // Account ID of the owner
	Capacity       int      `dynamodbav:"Capacity"`                 // 0 = no auto-confirm; >0 = total spots
}

func (d *Dto) ToDomain() matchoffer.Entity {
	// Convert timestamps to time.Time using the location's timezone
	var location matchoffer.Location
	if d.Latitude != nil && d.Longitude != nil {
		location = matchoffer.NewLocationWithCoords(d.Country, d.Province, d.Locality, *d.Latitude, *d.Longitude)
	} else {
		location = matchoffer.NewLocation(d.Country, d.Province, d.Locality)
	}
	tz := location.GetTimezone()

	day := time.Unix(d.Day, 0).In(tz)
	startTime := time.Unix(d.StartTime, 0).In(tz)
	endTime := time.Unix(d.EndTime, 0).In(tz)
	createdAt := time.Unix(d.CreatedAt, 0).In(tz)

	timeSlot, _ := matchoffer.NewTimeSlot(startTime, endTime)

	// Reconstruct CategoryRange based on RangeType
	var categoryRange matchoffer.CategoryRange
	switch matchoffer.RangeType(d.RangeType) {
	case matchoffer.RangeTypeSpecific:
		categories := make([]common.Category, len(d.Categories))
		for i, c := range d.Categories {
			categories[i] = common.Category(c)
		}
		categoryRange = matchoffer.NewSpecificCategories(categories)
	case matchoffer.RangeTypeGreaterThan:
		categoryRange = matchoffer.NewGreaterThanCategory(common.Category(d.MinLevel))
	case matchoffer.RangeTypeLessThan:
		categoryRange = matchoffer.NewLessThanCategory(common.Category(d.MaxLevel))
	case matchoffer.RangeTypeBetween:
		categoryRange, _ = matchoffer.NewBetweenCategories(common.Category(d.MinLevel), common.Category(d.MaxLevel))
	}

	status, _ := matchoffer.ParseStatus(d.Status)

	return matchoffer.Entity{
		ID:                 d.Id,
		TeamName:           d.TeamName,
		Sport:              common.Sport(d.Sport),
		Day:                day,
		TimeSlot:           timeSlot,
		Location:           location,
		AdmittedCategories: categoryRange,
		Status:             status,
		CreatedAt:          createdAt,
		OwnerAccountID:     d.OwnerAccountId,
		Capacity:           d.Capacity,
	}
}
