package matchannouncement

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	"time"
)

type Dto struct {
	EntityId   string `dynamodbav:"EntityId"`   // "Entity#MatchAnnouncement"
	Id         string `dynamodbav:"Id"`         // Generated UUID
	TeamName   string `dynamodbav:"TeamName"`   // Team name
	Sport      string `dynamodbav:"Sport"`      // Sport type
	Day        int64  `dynamodbav:"Day"`        // Unix timestamp of the day
	StartTime  int64  `dynamodbav:"StartTime"`  // Unix timestamp of start time
	EndTime    int64  `dynamodbav:"EndTime"`    // Unix timestamp of end time
	Country    string `dynamodbav:"Country"`    // Location country
	Province   string `dynamodbav:"Province"`   // Location province
	Locality   string `dynamodbav:"Locality"`   // Location locality
	RangeType  string `dynamodbav:"RangeType"`  // Category range type
	Categories []int  `dynamodbav:"Categories"` // Specific categories
	MinLevel   int    `dynamodbav:"MinLevel"`   // Minimum category level
	MaxLevel   int    `dynamodbav:"MaxLevel"`   // Maximum category level
	Status     string `dynamodbav:"Status"`     // Announcement status
	CreatedAt  int64  `dynamodbav:"CreatedAt"`  // Unix timestamp of creation
	ExpiresAt  int64  `dynamodbav:"ExpiresAt"`  // TTL for DynamoDB (Unix timestamp)
}

func (d *Dto) ToDomain() matchannouncement.Entity {
	// Convert timestamps to time.Time using the location's timezone
	location := matchannouncement.NewLocation(d.Country, d.Province, d.Locality)
	tz := location.GetTimezone()

	day := time.Unix(d.Day, 0).In(tz)
	startTime := time.Unix(d.StartTime, 0).In(tz)
	endTime := time.Unix(d.EndTime, 0).In(tz)
	createdAt := time.Unix(d.CreatedAt, 0).In(tz)

	timeSlot, _ := matchannouncement.NewTimeSlot(startTime, endTime)

	// Reconstruct CategoryRange based on RangeType
	var categoryRange matchannouncement.CategoryRange
	switch matchannouncement.RangeType(d.RangeType) {
	case matchannouncement.RangeTypeSpecific:
		categories := make([]common.Category, len(d.Categories))
		for i, c := range d.Categories {
			categories[i] = common.Category(c)
		}
		categoryRange = matchannouncement.NewSpecificCategories(categories)
	case matchannouncement.RangeTypeGreaterThan:
		categoryRange = matchannouncement.NewGreaterThanCategory(common.Category(d.MinLevel))
	case matchannouncement.RangeTypeLessThan:
		categoryRange = matchannouncement.NewLessThanCategory(common.Category(d.MaxLevel))
	case matchannouncement.RangeTypeBetween:
		categoryRange, _ = matchannouncement.NewBetweenCategories(common.Category(d.MinLevel), common.Category(d.MaxLevel))
	}

	status, _ := matchannouncement.ParseStatus(d.Status)

	return matchannouncement.Entity{
		ID:                 d.Id,
		TeamName:           d.TeamName,
		Sport:              common.Sport(d.Sport),
		Day:                day,
		TimeSlot:           timeSlot,
		Location:           location,
		AdmittedCategories: categoryRange,
		Status:             status,
		CreatedAt:          createdAt,
	}
}
